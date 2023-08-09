import random
import json
import sys
import os

if __name__ == "__main__":
    raw_file = sys.argv[1]
    if not raw_file or not os.path.exists(raw_file):
        print("Must specify raw geojson file as first argument")
        exit(1)
    clean_file_name = raw_file.replace("raw", "clean")
    raw_data_file = open(sys.argv[1])
    raw_data_json = raw_data_file.read()
    raw_data_file.close()
    raw_data_dict = json.loads(raw_data_json)
    print('features before sample', len(raw_data_dict['features']))
    samples = min(len(raw_data_dict['features']), 500)
    print(f'taking {samples} samples')
    raw_data_dict['features'] = random.sample(
        raw_data_dict['features'], samples)
    clean_file = open(clean_file_name, 'x')
    clean_data_json = json.dumps(raw_data_dict)
    clean_file.write(clean_data_json)
    clean_file.close()
