import requests
import os
from subprocess import call
from multiprocessing import Pool


def file_key(f):
    if not '-' in f:
        return -1
    else:
        return int(f.split('-')[1].split('.')[0])


def clean_dir(path):
    if os.path.exists(path):
        for filename in os.listdir(path):
            os.remove(path + filename)

    if not os.path.exists(path):
        os.mkdir(path)


def to_ocr(pdf_path):
    base_dir = '/tmp/pdf_to_ocr_out/'
    clean_dir(base_dir)

    assert os.path.exists(base_dir)
    assert len(os.listdir(base_dir)) == 0

    image_out = base_dir + 'out.png'
    command = 'convert -density 300 ' + pdf_path + ' -quality 100 ' + image_out
    call(command.split())
    print('image conversion done')

    with open(base_dir + 'png.list', 'w') as png_list:
        files = [f for f in os.listdir(base_dir) if f.endswith('.png')]
        files = sorted(files, key=file_key)
        png_list.write('\n'.join([base_dir + f for f in files]))

    command = 'tesseract ' + base_dir + 'png.list ' + base_dir + 'ocr'
    call(command.split())

    with open(base_dir + 'ocr.txt', 'r') as ocr:
        text = ocr.read()
    print('ocr done')

    clean_dir(base_dir)

    return text


def update_document_ocr(id, api_key,
                        url='http://allenarchive-dev.iac.gatech.edu/api/'):
    eqs = '===================================================================='
    print(eqs)
    print('Document', id)
    print(eqs)
    req = requests.get(url + 'files/' + str(id))
    if len(req.content) > 50000:
        return
    item = req.json()
    print('file retrieved')

    pdf_path = '/tmp/ocr_pdf_out.pdf'
    if 'file_urls' in item:
        pdf = requests.get(item['file_urls']['original'])
        with open(pdf_path, 'wb') as out:
            out.write(pdf.content)
            print('pdf written')

    ocr_text = to_ocr(pdf_path)

    req = requests.get(url + 'items/' + str(id))
    item = req.json()
    print('item retrieved')

    if not 'element_texts' in item:
        return
    for i, element_text in enumerate(item['element_texts']):
        if element_text['element']['name'] == 'Text' \
                and element_text['element_set']['name'] == 'Item Type Metadata':
            item['element_texts'][i]['text'] = ocr_text

    req = requests.put(url + 'items/' + str(id) + '?key=' + api_key, json=item)
    print('update request sent')


if __name__ == "__main__":
    for i in range(9700, 12000):
        update_document_ocr(i, 'aa516a5f41a594de03b8d9ed1552dc5847a6ac9a')
