<?php $title = __('Autotagging');?>
<?php echo head(array('title' => __('Autotagging'), 'bodyclass' => 'w3Page'));?>
<head>
    <script>
        window.onload = function() {
            var autotag = document.getElementById('autotag'); 
            autotag.onclick = autoTagDocuments;
        };

        var autoTagDocuments = function() {
            var details = document.getElementById('details');
            var start = document.getElementById('start').value;
            var end = document.getElementById('end').value;
            for (var i = start; i <= end; i++) {
                var data = {
                    'action': 'tag',
                    'start': i,
                    'end': i + 1,
                    'url': "<?php echo rtrim(absolute_url(""), "admin/") ?>",
                    'api_key': document.getElementById('api_key').value
                };
                var status = document.getElementById('status');
                status.innerHTML = '<strong>Status: In progress.</strong>';
                status.style.display = 'block';
                jQuery.post('/admin/autotagging/index/autotag', data, function (response) {
                    details.innerHTML += response;
                    status.innerHTML = '<strong>Status: Done tagging.';
                    var expander = document.getElementById('details_expander');
                    expander.style.display = 'block';
                    expander.onclick = function() {
                        if (details.style.display === 'none') {
                            details.style.display = 'block';
                        } else {
                            details.style.display = 'none';
                        }
                    };
               });
            }
        };
    </script>
</head>
<body>
    <p>
        Enter your tagging API key below. 
        Then enter a range of documents to be tagged, e.g., documents 1 through 10.
        The tagging process may take a while.
    </p>
    <input type="text" placeholder="Tagging API key" id="api_key"/>
    <input type="number" placeholder="Start ID" id="start"/>
    <input type="number" placeholder="End ID" id="end"/>
    <button id="autotag">Tag Documents</button> 
    <p id="status" style="display: none;"></p>
    <button id="details_expander" style="display: none;">Display Details</button>
    <pre id="details" style="display: none;"></pre>
</body>

