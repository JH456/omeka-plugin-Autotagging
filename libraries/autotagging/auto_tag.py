import argparse
import json

import nltk
import requests
import spacy

nlp = spacy.load('en_core_web_sm')


def get_labeled_entities_nltk(document):
    sentences = nltk.sent_tokenize(document)
    sentences = [nltk.pos_tag(nltk.word_tokenize(sent)) for sent in sentences]

    labels = []
    entities = []
    for sent in sentences:
        chunk = nltk.ne_chunk(sent)
        labels.extend(
            tree.label() for tree in chunk if isinstance(tree, nltk.Tree)
        )
        entities.extend(
            ' '.join(i[0] for i in tree.leaves()).lower()
            for tree in chunk if isinstance(tree, nltk.Tree)
        )

    return set(zip(labels, entities))


def get_labeled_entities_spacy(document):
    text = nlp(document)
    labels = []
    entities = []
    for e in set(t for t in text.ents if str(t).strip()):
        entities.append(' '.join(str(e).lower().split()))
        labels.append(e.label_)
    return set(zip(labels, entities))


def tag_document(id, api_key, url, label_mapping,
                 get_labeled_entities=get_labeled_entities_spacy,
                 max_length=50000):
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

    if 'element_texts' in item:
        element_texts = item['element_texts']
        texts = [i['text'] for i in element_texts
                 if i['element']['name'] == 'Text']
        if texts and len(texts[0]) < max_length:
            labeled_entities = get_labeled_entities(texts[0])
            tags = [label_mapping[label] + ': ' + entity
                    for label, entity in labeled_entities
                    if label in label_mapping]

            for tag in tags:
                item['tags'].append({'resource': 'tags', 'name': tag})

    print('tags:', json.dumps(item['tags'], sort_keys=True, indent=4))
    req = requests.put(url + 'items/' + str(id) + '?key=' + api_key, json=item)
    print('update request sent')

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
    parser.add_argument('-l', '--max_length',
                        help='Maximum length in characters of documents',
                        type=int, default=50000)
    parser.add_argument('-t', '--tag_engine',
                        help='NLP engine for tagging. nltk for NLTK, spacy for spaCy',
                        type=str, default='spacy')

    args = parser.parse_args()
    api_key = args.key

    if args.tag_engine == 'nltk':
        get_labeled_entities = get_labeled_entities_nltk
    else:
        get_labeled_entities = get_labeled_entities_spacy

    label_mapping = {
        'PERSON': 'Person',
        'FACILITY': 'Facility',
        'DATE': 'Date',
        'ORG': 'Organization',
        'ORGANIZATION': 'Organization',
        'GPE': 'Geopolitical Entity',
        'LCO': 'Location',
        'EVENT': 'Event',
        'LAW': 'Law'
    }
    for i in range(args.start, args.end):
        tag_document(i, args.key, args.url, label_mapping,
                     max_length=args.max_length or 50000,
                     get_labeled_entities=get_labeled_entities)
