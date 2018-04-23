import argparse

import requests
import spacy

nlp = spacy.load('en_core_web_sm')


def tag_document(id, api_key, url):
    eqs = '===================================================================='
    print(eqs)
    print('Document', id)
    print(eqs)
    url = url.rstrip('/') + '/api/'
    req = requests.get(url + 'items/' + str(id))
    item = req.json()
    if 'message' in item and item['message'].endswith('Record not found.'):
        print('Record with id ' + str(id) + ' could not be found.')
        return
    entity_mapping = {
        'PERSON': 'Person',
        'FACILITY': 'Facility',
        'DATE': 'Date',
        'ORG': 'Organization',
        'GPE': 'Geopolitical Entity',
        'LCO': 'Location',
        'EVENT': 'Event',
        'LAW': 'Law'
    }
    if 'element_texts' in item:
        element_texts = item['element_texts']
        texts = [i['text'] for i in element_texts
                 if i['element']['name'] == 'Text']
        if texts and len(texts[0]) < 50000:
            tags = []
            text = nlp(texts[0])
            for e in set(t for t in text.ents if str(t).strip()):
                if e.label_ in entity_mapping:
                    label = entity_mapping[e.label_]
                    tags.append(label + ': ' +
                                ' '.join(str(e).lower().split()))

            for tag in tags:
                item['tags'].append({'resource': 'tags', 'name': tag})

    req = requests.put(url + 'items/' + str(id) + '?key=' + api_key, json=item)


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description='Generate NER tags on a set of Omeka items.'
    )
    parser.add_argument('url', help='Base URL for the Omeka instance')
    parser.add_argument('key', help='API key for the Omeka instance')
    parser.add_argument('-s', '--start',
                        help='Document ID where tagging should begin',
                        type=int, default=0)
    parser.add_argument('-e', '--end',
                        help='Document ID where tagging should end',
                        type=int, default=100000)

    args = parser.parse_args()
    api_key = args.key

    for i in range(args.start, args.end):
        tag_document(i, args.key, args.url)
