<?php
    $title = __('Autotagging');
    // bodyclass is a custom class for the body
    function debug_to_console( $data ) {
        $output = $data;
        if ( is_array( $output ) )
            $output = implode( ',', $output);
    
        echo "<script>console.log( 'Debug Objects: " . $output . "' );</script>";
    }
?>

<pre>
<?php 
    $text = "This is a test of the georgia tech emergency\nnotification system. I repeat; this is only\na test. don't worry too much.";
    $text = str_replace("\n", " ", $text);
    $wrapped = wordwrap($text, 100);
    echo explode("\n", $wrapped)[0];
    ?>
</pre>

<p1><?php $command ='python3 ./plugins/Autotagging/views/public/pythons/data_gatherer.py';
    exec($command, $out, $status);
    foreach ($out as $value) {
        debug_to_console($value);
    }
?></p1>

<h1><?php echo $title; ?></h1>


<strong>Coming soon! Tagging works in python tho!</strong>


