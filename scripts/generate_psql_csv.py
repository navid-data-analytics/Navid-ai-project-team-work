from optparse import OptionParser

import binascii
from postgres import Postgres
from getpass import getpass
import pandas
import time
import csv
import hvac
import os
from ctypes import string_at
from sys import getsizeof
from binascii import hexlify

TIMESTAMP = time.time()
pandas.set_option('display.max_colwidth', -1)


def read_creds(vault_username, vault_password, db_type):
    client = hvac.Client('https://c1-vault.aws.callstats.io')
    client.auth_userpass(
        vault_username, vault_password, mount_point='ops/userpass')
    if db_type == 'test':
        path = 'test/postgresql/aid_service/creds/aid_service_readonly'
    elif db_type == 'prod':
        path = 'prod/postgresql/aid_service/creds/aid_service_readonly'
    else:
        raise AssertionError('Unknown database type, not in ("prod", "test")!')
    creds = client.read(path)

    username = creds['data']['username']
    password = creds['data']['password']
    return (username, password)


parser = OptionParser(description='Get psql data')

parser.add_option('--username', type=str)
parser.add_option('--password', type=str)
parser.add_option('--db_type', type=str)
args, _ = parser.parse_args()

vault_username = args.username
vault_password = args.password
db_type = args.db_type

if args.username is None:
    vault_username = input('Username: ')

if args.password is None:
    vault_password = getpass('Password: (will be hidden) ')

if args.db_type is None:
    db_type = input('db_type (test or prod): ')

db_username, db_password = read_creds(vault_username, vault_password, db_type)

print("Using {} {}".format(db_username, db_password))

db_uri = "postgres://{}:{}@aid-{}.caaddmbzjhml.eu-west-1.rds.amazonaws.com:5432/aid".format(  # noqa
    db_username, db_password, db_type)

print("Connecting to psql database {}".format(db_uri))
db = Postgres("{}".format(db_uri))
print("Fetching data...")
out = db.all(
    "SELECT data, id, app_id, template_id, generated_at FROM messages;")
template_id_mapping = db.all(
    "SELECT messages.template_id, type FROM messages INNER JOIN message_templates ON (messages.template_id = message_templates.id);"  # noqa
)

template_id_mapping = {
    tim.template_id: tim.type
    for tim in template_id_mapping
}

print('Parsing psql data...')

# Parse IDs
ids = [message.id for message in out]

# Parse appIDs
appids = [message.app_id for message in out]

# Parse templateIDs
templateids = [template_id_mapping[message.template_id] for message in out]

# Parse generated at
generatedat = [message.generated_at for message in out]

# Parse data
data = [str(message.data.tobytes()).replace(',', ' ') for message in out]

assert len(ids) == len(data) == len(appids) == len(templateids) == len(
    generatedat), 'Not all datalen equal! i{} d{} a{} t{} g{}'.format(
        len(ids), len(data), len(appids), len(templateids), len(generatedat))

print('Writing data to HTML...')
csv_path = '{}_psql_data_{}.csv'.format(db_type, TIMESTAMP)
try:
    with open(csv_path, 'w') as csvfile:
        csvfile.write('ID,appID,templateID,generatedAt,data\n')
        for id, appid, template, gen, dat in zip(ids, appids, templateids,
                                                 generatedat, data):
            csvfile.write('{},{},{},{},{}\n'.format(id, appid, template, gen,
                                                    dat))
        df = pandas.read_csv(csv_path)
        print('Save data as html table...')
        df.to_html('{}_psql_data_{}.html'.format(db_type, TIMESTAMP))
except BaseException as e:
    print('Something went wrong:')
    print(e)
    print('Csv file corrupted, removing...')
finally:
    os.remove(csv_path)
