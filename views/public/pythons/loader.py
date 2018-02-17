import os

def load_data():
	# data_dir = '/home/kpberry/Desktop/iada_pdfs/originals/'
	data_dir = '/home/kpberry/Desktop/iada_pdfs/test/'
	print('start')
	documents = []
	for file in os.listdir(data_dir):
		if file.endswith('.txt'):
			with open(data_dir + file, 'r') as in_file:
				lines = ''.join([line for line in in_file])
				documents.append(lines)
	return documents