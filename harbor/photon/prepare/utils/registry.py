import os, copy

from g import config_dir, templates_dir, DEFAULT_GID, DEFAULT_UID, data_dir
from utils.misc import prepare_dir
from utils.jinja import render_jinja


registry_config_dir = os.path.join(config_dir, "registry")
registry_config_template_path = os.path.join(templates_dir, "registry", "config.yml.jinja")
registry_conf = os.path.join(config_dir, "registry", "config.yml")
registry_data_dir = os.path.join(data_dir, 'registry')

levels_map = {
    'debug': 'debug',
    'info': 'info',
    'warning': 'warn',
    'error': 'error',
    'fatal': 'fatal'
}

def prepare_registry(config_dict):
    prepare_dir(registry_data_dir, uid=DEFAULT_UID, gid=DEFAULT_GID)
    prepare_dir(registry_config_dir)

    storage_provider_info = get_storage_provider_info(
    config_dict['storage_provider_name'],
    config_dict['storage_provider_config'])

    render_jinja(
        registry_config_template_path,
        registry_conf,
        uid=DEFAULT_UID,
        gid=DEFAULT_GID,
        level=levels_map[config_dict['log_level']],
        storage_provider_info=storage_provider_info,
        **config_dict)


def get_storage_provider_info(provider_name, provider_config):
    provider_config_copy = copy.deepcopy(provider_config)
    if provider_name == "filesystem":
        if not (provider_config_copy and ('rootdirectory' in provider_config_copy)):
            provider_config_copy['rootdirectory'] = '/storage'
    if provider_name == 'gcs' and provider_config_copy.get('keyfile'):
        provider_config_copy['keyfile'] = '/etc/registry/gcs.key'
    # generate storage configuration section in yaml format
    storage_provider_conf_list = [provider_name + ':']
    for config in provider_config_copy.items():
        if config[1] is None:
            value = ''
        elif config[1] == True:
            value = 'true'
        else:
            value = config[1]
        storage_provider_conf_list.append('{}: {}'.format(config[0], value))
    storage_provider_info = ('\n' + ' ' * 4).join(storage_provider_conf_list)
    return storage_provider_info
