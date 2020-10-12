import shutil, os

from g import config_dir, templates_dir, data_dir, DEFAULT_GID, DEFAULT_UID
from utils.misc import prepare_dir, generate_random_string
from utils.jinja import render_jinja

core_config_dir = os.path.join(config_dir, "core", "certificates")
core_env_template_path = os.path.join(templates_dir, "core", "env.jinja")
core_conf_env = os.path.join(config_dir, "core", "env")
core_conf_template_path = os.path.join(templates_dir, "core", "app.conf.jinja")
core_conf = os.path.join(config_dir, "core", "app.conf")

ca_download_dir = os.path.join(data_dir, 'ca_download')
psc_dir = os.path.join(data_dir, 'psc')


def prepare_core(config_dict, with_notary, with_clair, with_chartmuseum):
    prepare_dir(psc_dir, uid=DEFAULT_UID, gid=DEFAULT_GID)
    prepare_dir(ca_download_dir, uid=DEFAULT_UID, gid=DEFAULT_GID)
    prepare_dir(core_config_dir)
    # Render Core
    # set cache for chart repo server
    # default set 'memory' mode, if redis is configured then set to 'redis'
    if len(config_dict['redis_host']) > 0:
        chart_cache_driver = "redis"
    else:
        chart_cache_driver = "memory"

    render_jinja(
        core_env_template_path,
        core_conf_env,
        chart_cache_driver=chart_cache_driver,
        with_notary=with_notary,
        with_clair=with_clair,
        with_chartmuseum=with_chartmuseum,
        **config_dict)

    render_jinja(
        core_conf_template_path,
        core_conf,
        uid=DEFAULT_UID,
        gid=DEFAULT_GID,
        xsrf_key=generate_random_string(40))



def copy_core_config(core_templates_path, core_config_path):
    shutil.copyfile(core_templates_path, core_config_path)
    print("Generated configuration file: %s" % core_config_path)
