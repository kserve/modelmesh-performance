"""This file contains the constants related to K6 and Howitzer parameters. The
large dictionary for K6 options is especially important because it affects
how the Javascript will render. Be careful when making changes, and make sure
that all unit tests still pass if you make edits to this file.
"""

# List of valid K6 options. Options that are set to true may be passed in as
# iterables and exploded out. For most options, it doesn't make sense to add
# iterability - these options are set to False, and will not be exploded out.
K6_VALID_OPTS = {
    'duration' : {'iterable':True, 'type':str},
    'rps' : {'iterable':True, 'type':str},
    'iterations' : {'iterable':True, 'type':str},
    'vus' : {'iterable':True, 'type':str},

    'vusMax' : {'iterable':False, 'type':str},
    'ext' : {'iterable':False, 'type':dict},
    'hosts' : {'iterable':False, 'type':dict},
    'insecureSkipTLSVerify' : {'iterable':False, 'type':bool},
    'linger' : {'iterable':False, 'type':bool},
    'maxRedirects' : {'iterable':False, 'type':str},
    'batch' : {'iterable':False, 'type':str},
    'batchPerHost' : {'iterable':False, 'type':str},
    'noConnectionReuse' : {'iterable':False, 'type':bool},
    'noVUConnectionReuse' : {'iterable':False, 'type':bool},
    'noUsageReport' : {'iterable':False, 'type':bool},
    'paused' : {'iterable':False, 'type':bool},
    'stages' : {'iterable':False, 'type':list},
    'tags' : {'iterable':False, 'type':dict},
    'thresholds' : {'iterable':False, 'type':dict},
    'throw' : {'iterable':False, 'type':bool},
    'blacklistIPs' : {'iterable':False, 'type':list},
    'summaryTrendStats' : {'iterable':False, 'type':list},
    'tisAuth' : {'iterable':False, 'type':list},
    'tisCipherSuites' : {'iterable':False, 'type':list},
    'tisVersion' : {'iterable':False, 'type':dict},
    'userAgent' : {'iterable':False, 'type':str},
    'httpDebug' : {'iterable':False, 'type':str},
    'systemTags' : {'iterable':False, 'type':list},
    'setupTimeout' : {'iterable':False, 'type':str},
    'teardownTimeout' : {'iterable':False, 'type':str}
}

# Parameters corresponding to each Howtizer test in the config options.
HOWITZER_PARAMS = {
    'required': ['name', 'template', 'base_url', 'k6_opts', 'model_name'],
    'optional': ['payload', 'description']
}