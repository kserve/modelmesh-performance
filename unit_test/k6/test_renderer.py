import json
import os
import unittest
import shutil

from perf_test.scripts import renderer

def extract_value(response):
    val = response.split('=')[-1]
    return ''.join(val.split())

class TestActualRendering(unittest.TestCase):

    def setUp(self) -> None:
        """This method runs before each test begins to run"""
        self.render_dir = "tmp/renders"
        os.makedirs(self.render_dir, exist_ok=True)

    def tearDown(self) -> None:
        shutil.rmtree("tmp")
        pass

    def test_renderer_makes_the_correct_number_of_tests(self):
        # sample test config should result in 4 tests
        renderer.generate_k6_tests(self.render_dir, "unit_test/fixtures", "unit_test/fixtures/test_config.json")
        self.assertEqual(4, len(os.listdir(self.render_dir)))

    def test_the_exploded_test_configs_have_correct_parameters(self):
        renderer.generate_k6_tests(self.render_dir, "unit_test/fixtures", "unit_test/fixtures/test_config.json")
        renderer.write_test_configs_to_file("tmp/test_final_config.json")

        final_config = json.load(open("tmp/test_final_config.json", "r"))
        vus = [foo["test_configs"]["k6_opts"]["vus"] for foo in final_config.values()]
        self.assertEqual([5, 10, 25, 50], vus)

# Checks to ensure that all object types can be rendered into Node as params.
class DictionaryToNodeConversionTests(unittest.TestCase):
    def test_node_base_case(self):
        expected = '{};'
        extracted_response = extract_value(renderer.dict_to_node({}))
        self.assertEqual(extracted_response, expected)

    def test_node_string_value(self):
        expected = "{'duration':'120m'};"
        extracted_response = extract_value(renderer.dict_to_node({'duration':'120m'}))
        self.assertEqual(extracted_response, expected)

    def test_node_bool_value(self):
        expected = "{'insecureSkipTLSVerify':true};"
        extracted_response = extract_value(renderer.dict_to_node({'insecureSkipTLSVerify': True}))
        self.assertEqual(extracted_response, expected)

    def test_node_list_value(self):
        input = {'stages':[{'duration':'1m','target':10},{'duration':'1m','target':0}]}
        expected ="{'stages':[{'duration':'1m','target':10},{'duration':'1m','target':0}]};"
        extracted_response = extract_value(renderer.dict_to_node(input))
        self.assertEqual(extracted_response, expected)

    def test_node_dictionary_value(self):
        input = {'thresholds': {'http_req_duration': ['avg<100', 'p(95)<200']}}
        expected = "{'thresholds':{'http_req_duration':['avg<100','p(95)<200']}};"
        extracted_response = extract_value(renderer.dict_to_node(input))
        self.assertEqual(extracted_response, expected)

# Flattening is used to consolidate subtests. Checks arbitrary nesting cases.
class listFlatteningTests(unittest.TestCase):
    def test_flatten_empty_list(self):
        self.assertEqual([], renderer.flatten([]))

    def test_flatten_nested_empty_lists(self):
        self.assertEqual([], renderer.flatten([[[],[]], [[]]]))

    def test_flatten_nested_nonempty_lists(self):
        self.assertEqual([1,2,3,4], renderer.flatten([[1, 2, [],[3]], [[], 4]]))

# Checks basic validation for config key structure. Aside from the k6 options,
# the validity of other option values is not considered here.
class configKeyValidationTests(unittest.TestCase):
    def test_invalid_config_param(self):
        self.assertRaises(KeyError, renderer.validate_tests, [{'garbage':'_'}])

    def test_empty_config(self):
        self.assertRaises(KeyError, renderer.validate_tests, [{
            'name':'test',
            'template': 'a_template',
            'base_url': 'a_test_url',
            'k6_opts':{'vus':10},
            'invalid': '_'
        }])

    def test_valid_config(self):
        renderer.validate_tests({'test_configs':[{
            'name':'test',
            'template': 'a_template',
            'base_url': 'a_test_url',
            'k6_opts':{'vus':10},
            'model_name': 'a_model'
        }]})

    def test_invalid_k6_config(self):
        self.assertRaises(KeyError, renderer.validate_tests, [{
            'name':'test',
            'template': 'a_template',
            'base_url': 'a_test_url',
            'k6_opts':{'garbage': '_'},
        }])

if __name__ == '__main__':
    unittest.main()
