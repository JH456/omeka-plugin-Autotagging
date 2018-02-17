import re

from loader import load_data

date_pattern = re.compile(
	r'((?:0?[1-9]|1[012]|Jan[uary]*|Feb[ruay]*|Mar[ch]*|Apr[il]*|May|June*'
	r'|July*|Aug[ust]*|Sep[tembr]*|Oct[ober]*|Nov[embr]*|Dec[embr]*)'
	r'[- /.,]+(?:0[1-9]|[12][0-9]|3[01])'
	r'[- /.,]+(?:\d*))', re.I)

address_pattern = re.compile(r'((?:(?:\d+)\s+(?:[A-Z][a-z.]+(?:,| +)*?)+)' \
                             + '|(?:PO Box \d+))?(?:,|\s)*' \
                             + '((?:[A-Z][a-z. ]+)+)?(?:,|\s)*(\w\w)(?:,|\s)*?(\d{5})')

full_line_name_pattern = re.compile(
	r'^\s*'
	r'(\s*(?:Mr|Mrs|Ms|Dr|Dear)\.?)*'  # title
	r'(\s*[A-Z][a-z]*)\s*?'  # first name
	r'((?:\s+[A-Z]\.?|)*)'  # middle initials # need space if no period
	r'(\s*[A-Z][a-z]*)?'  # last name
	r'\s*$'
)


def find_locations(document):
	return address_pattern.findall(document)


def find_dates(document):
	return date_pattern.findall(document)


def find_full_line_names(document):
	names = []
	for line in document.split('\n'):
		if len(line.strip()) > 0:
			match = full_line_name_pattern.match(line)
			if match:
				names.append(match.groups())
	return names


def run_test():
	print(' loading data ...')
	documents = load_data()

	# results = [find_full_line_names(d) for d in documents]
	# results = [find_dates(d) for d in documents]
	# results = [find_locations(d) for d in documents]
	print(len(documents))

	# print(len(results))


if __name__ == '__main__':
	print('... starting ...')
	run_test()
