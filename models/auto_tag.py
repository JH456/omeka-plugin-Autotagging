import spacy
import requests
from multiprocessing import Pool


nlp = spacy.load('en_core_web_sm')


def tag_document(id, api_key,
                 url='http://allenarchive-dev.iac.gatech.edu/api/items/'):
    eqs = '===================================================================='
    print(eqs)
    print('Document', id)
    print(eqs)
    req = requests.get(url + str(id))
    item = req.json()
    entity_mapping = {
        'PERSON': 'Person',
        'FACILITY': 'Facility',
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

    req = requests.put(url + str(id) + '?key=' + api_key, json=item)


def td(d):
    tag_document(d, 'aa516a5f41a594de03b8d9ed1552dc5847a6ac9a')


if __name__ == "__main__":
    pool = Pool(processes=32)
    pool.map(td, range(0, 12000))
