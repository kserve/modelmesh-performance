"""This module explodes the JSON configuration file into a flattened list
of K6 tests to be run. It then creates a rendering directory and renders
the K6 options into Javascript. Beware that this file must currently be
run using Python3; because it uses the json module to load the initial
dictionary, Python2 will load the strings as unicode, which will create
invalid Javascript when converted into a string (due to nested u'XXX' vals).
"""
from itertools import product
import os
import json
import shutil
from jinja2 import Environment, FileSystemLoader
from perf_test.scripts import constants
import argparse
import copy

TEST_CONFIGS_FILENAME = 'test_configs.dict'
TEST_CONFIGS = {}

def load_configs(config_file: str):
    """Loads the configurations test & explodes iterable parameters
    into a list of all combinations of test parameters. Returns a list
    of subtests - all of the K6 tests that need to be run to get results for
    all possible combinations of iterable parameters.

    Returns:
         list: The list of K6 options to be rendered into K6 template files.
    """
    configs = json.load(open(config_file, 'r'))
    validate_tests(configs)

    subtests = flatten([explode_test(test) for test in configs['test_configs']])
    return subtests

def validate_tests(configs):
    """Simple validation for the Howitzer configuration file, which checks the
    following for every config level test (before exploding to subtests).
        - The config tests do not have any invalid parameters.
        - The config tests have all required fields defined (i.e. name, template,
          base_url, and k6_opts).
        - All of the K6_opts map to a key in the dictionary
          defining all of the current K6 parameters.

    Args:
        configs (Dict): Config file for loading K6 test and rendering documents.
        required_fields (dict): Dictionary whose keys are ALL valid k6 opts.
    """
    if 'test_configs' not in configs:
        raise KeyError('test_configs key is missing from the configs!')
    for test in configs['test_configs']:
        if not set(constants.HOWITZER_PARAMS['required']).issubset(set([k for k in test.keys()])):
            raise KeyError('Missing required test descriptor(s) in test: ' + str(test))
        if not all(opt in constants.K6_VALID_OPTS.keys() for opt in test['k6_opts']):
            raise KeyError('One of the following K6 params is unknown: '+ str(test['k6_opts']))

def explode_test(test):
    """Blows up a single test configuration containing any number of valid iterable
    params. External parameters, which are the params not directly rendered into
    a K6 template (i.e. the name of the template) are duplicated across all resulting
    subtests when a iterable parameter is exploded into multiple K6 cofigurations.
    A toy example of how the explosion process produces multiple K6 tests for
    multiple iterable params is shown below.

    Args:
        test (dict): dictionary retrieved from JSON corresponding to one Howitzer test.

    Returns:
        list of options to be run in K6 needed to meet the iterable test config.

    Example:
        -----Config contents-----      -----resulting K6 options-----
        vus: [10, 100]                  vus: 10, duration: 120s
        duration: [120s, 300s]          vus: 10, duration: 300s
                                        vus: 100, duration: 120s
                                        vus: 100, duration: 300s
    """
    # Extract the external parameters, the static (i.e. already valid in K6) params
    # and the dynamic (i.e. iterable params that need to be split directly) params.
    external = {k: v for (k, v) in test.items() if k != 'k6_opts'}
    static = {k: v for (k, v) in test['k6_opts'].items() if not is_dynamic_param(k, v)}
    dynamic_raw = {k: v for (k, v) in test['k6_opts'].items() if is_dynamic_param(k, v)}
    # Take all possible combinations of dynamic parameters values, and use order
    # preservation to map the combination generator to its field name.
    dynamic_combos = product(*[test['k6_opts'][k] for k in [x for x in dynamic_raw]])
    dynamics = [{k:v for (k, v) in zip(dynamic_raw, val)} for val in dynamic_combos]
    # Combine external, static, and dynamic values to form the list of k6 tests
    return [get_static_subtest(static, dynamic, external) for dynamic in dynamics]

def flatten(lyst):
    """Recursively flattens an arbitrarily nested list.

    Args:
        lyst (list): List of arbitrarily nested objects.

    Returns:
        fully flattened list object.
    """
    if not isinstance(lyst, list):
        return [lyst]
    elif not lyst:
        return []
    elif len(lyst) > 1:
        return flatten(lyst[0]) + flatten(lyst[1:])
    return flatten(lyst[0]) + []

def get_static_subtest(static_fields, dynamic_params, external_params):
    """Combines iterable parameters for a single K6 test (post-explosion) with static
    parameters, thereby encompassing a complete description for a K6 test. Creates
    a dictionary with the k6 test information, and test info external to K6 (such
    as the template which the test will be rendered into).

    Args:
        subtest (tuple): raw values corresponding to iterable/dynamic k6 params.
        static_fields (dict): static test params (valid k6 opts without exploding).
        dynamic_params (dict): dynamic test params (iterable vals that explode into k6 params).

    Returns:
        dictionary object encapsulating all object relevant to a single load test.
    """
    extracted_subtest = dict(dynamic_params, **static_fields)
    return {'k6_test':extracted_subtest, 'external_params':external_params}

def is_dynamic_param(key, val):
    """Checks whether or not a retrieved load testing parameter is iterable
    (i.e. configured as a list in the config JSON file, and marked as iterable
    in the constants map to K6 parameters). If a load testing parameter is marked
    as iterable, then it will explode into multiple load tests.

    Args:
        key (dict): k6 option being considered.
        val (any): The value that the key maps to in the given configs.

    Returns:
        bool: True if the key allows dynamic opts and the value is a list.
    """
    try:
        return constants.K6_VALID_OPTS[key]['iterable'] and isinstance(val, list)
    except KeyError:
        raise KeyError('Could not check if key: ' + str(key) + ' was iterable. Is it a K6 opt?')


def sweep_render_dir(render_dir: str):
    """Checks to see if the rendering directory exists. If it does, it removes
    it and recreates a new rendering directory.
    """
    if os.path.exists(render_dir):
        shutil.rmtree(render_dir)
    os.mkdir(render_dir)


def update_test_configs(test_name, info):
    """Updates the global TEST_CONFIGS for each test case scenario
    This is used to generate the JSON summary report

    Args:
        test_name (string): test case name with the scenario number
        info (dict): test parameters after Howitzer config explosion containing k6_test
         and external_params
    """
    TEST_CONFIGS[test_name] = {}
    # We need to deep copy here because across tests, info['external_params']
    # remains the same. Without the deepcopy, we would make each
    # TEST_CONFIGS[test_name]['test_configs'] reference the same object
    # which would make any further changes to `k6_opts` affect all test configs
    TEST_CONFIGS[test_name]['test_configs'] = copy.deepcopy(info['external_params'])
    TEST_CONFIGS[test_name]['test_configs']['k6_opts'] = info['k6_test']

def render_k6_test(info, render_dir: str, template_dir: str):
    """Renders the k6 Javascript options into the test's corresponding template.
    The runnable k6 test can be found in the rendered directory.

    Args:
        info (dict): test parameters after Howitzer config explosion containing k6_test
         and external_params
    """
    params = info['external_params']
    k6_opts = dict_to_node(info['k6_test'])
    test_name = params['name']
    env = Environment(loader=FileSystemLoader(template_dir))
    template = env.get_template(params['template'])
    rendered_test = template.render(params, k6_opts=k6_opts)
    scenario_num = len([x for x in os.listdir(render_dir) if x.startswith(test_name)])
    filename = os.path.join(render_dir, test_name + '_' + str(scenario_num) + '.js')
    with open(filename, 'w') as render_file:
        render_file.write(rendered_test)
    update_test_configs(test_name + '_' + str(scenario_num), info)


def dict_to_node(test):
    """Converts the options dictionary corresponding to a single k6 test into
    a javascript object that can be rendered directly into a template.

    Args:
        test (dict): k6 load test parameters to be converted to a renderable form.

    Returns:
        str: Javascript (K6 options object) as a string.
    """
    ret_str = 'export let options = {\n'
    for idx, k in enumerate(test):
        ret_str += render_js_value(k, test[k])
        if idx != len(test)-1:
            ret_str += ',\n'
    ret_str += '\n};'
    return ret_str

def render_js_value(key, val):
    """Renders a key value pair into a form that Javascript can parse it.
    Formatting is dependent on the type mapping for each parameter, as
    specified in the constants file.

    Args:
        key (str): name of the k6 parameters being rendered into Javascript.
        val (any): value of the k6 parameter being rendered into Javascript.

    Returns:
        str: One key/val pair inside of the K6 opts JS object, represented as a string.
    """
    rendered_key = '    \'' + str(key) + '\': '
    if constants.K6_VALID_OPTS[key]['type'] in [int, str]:
        return  rendered_key + '\''+ str(val) + '\''
    elif constants.K6_VALID_OPTS[key]['type'] is bool:
        return rendered_key + str(val).lower()
    elif constants.K6_VALID_OPTS[key]['type'] in [list, dict]:
        return rendered_key + str(val)
    else:
        raise TypeError('Received an invalidly typed config. Could not render!')


def write_test_configs_to_file(test_configs_filename: str):
    """Write the test_configs corresponding to each rendered test in TEST_CONFIGS_FILENAME

    Args:
        None
    """
    with open(test_configs_filename, 'w') as config_file:
        json.dump(TEST_CONFIGS, config_file, indent=4)

def generate_k6_tests(render_dir: str, template_dir: str, config_file: str):
    """Primary function to be invoked for end to end rendering. Loads a parameter
    file (meant to be set via Kubernetes configmap), explodes the configs into
    a list of k6 tests, converts the resulting objects to Javascript, and renders
    them into the appropriate templates.
    """
    subtests = load_configs(config_file)
    sweep_render_dir(render_dir)
    for info in subtests:
        render_k6_test(info, render_dir, template_dir)

if __name__ == '__main__':

    '''Obtain metrics from sysdig and create reports.
    '''
    parser = argparse.ArgumentParser(add_help=False,
                                     description=
                                     'Render some mother trucking K6 tests.')

    parser.add_argument('--help', action='help', help='Show this help message and exit.')

    parser.add_argument('-r', '--render-dir', metavar='RENDER_DIR',
                        help='The directory to render k6 tests into',
                        required=True)
    parser.add_argument('-t', '--template-dir', metavar='TEMPLATE_DIR',
                        help='The directory to find the k6 test template in',
                        required=True)
    parser.add_argument('-s', '--summaries-dir', metavar='SUMMARIES_DIR',
                        help='The directory to place the k6 per-test summary files into',
                        required=True)
    parser.add_argument('-c', '--config-file', metavar='CONFIG_FILE',
                        help='The absolute path to the one true test config file to be exploded',
                        required=True)

    args = parser.parse_args()

    generate_k6_tests(args.render_dir, args.template_dir, args.config_file)

    # We slap a file that contains all the test configs (configs only, _not_ rendered tests) into the summaries
    # directory. The scraper will slam back together each test config with its k6 summary.
    test_configs_file = os.path.join(args.summaries_dir, TEST_CONFIGS_FILENAME)
    write_test_configs_to_file(test_configs_file)