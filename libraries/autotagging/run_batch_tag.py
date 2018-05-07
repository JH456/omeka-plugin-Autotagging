from subprocess import Popen

batches = [i for i in range(0, 2136, 100)]
batches = zip(batches[:-1], batches[1:])
for start, end in batches:
    Popen([
        'python3', 'auto_tag.py',
        'http://allenarchive-dev.iac.gatech.edu/',
        '030c516f3f818bb10793ff6c965489c69647129d',
        '-s', str(start),
        '-e', str(end)
    ])
