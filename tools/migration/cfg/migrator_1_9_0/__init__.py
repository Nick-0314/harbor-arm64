from __future__ import print_function
import utils
import os
import yaml
from jinja2 import Environment, FileSystemLoader, StrictUndefined

acceptable_versions = ['1.8.0']

def migrate(input_cfg, output_cfg):
    config_dict = utils.read_conf(input_cfg)

    this_dir = os.path.dirname(__file__)
    tpl = Environment(
        loader=FileSystemLoader(this_dir),
        undefined=StrictUndefined,
        trim_blocks=True,
        lstrip_blocks=True
        ).get_template('harbor.yml.jinja')

    with open(output_cfg, 'w') as f:
        f.write(tpl.render(**config_dict))