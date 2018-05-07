<?php
/**
 * Autotagging
 *
 * @license http://www.gnu.org/licenses/gpl-3.0.txt GNU GPLv3
 */

/**
 * The Autotagging index controller class.
 *
 * @package Autotagging
 */
class Autotagging_IndexController extends Omeka_Controller_AbstractActionController
{
    public function browseAction()
    {}

    public function autotagAction()
    {
        $start = 0;
        $end = 0;

        if (isset($_POST['start'])) {
            $start = $_POST['start'];
        }

        if (isset($_POST['end'])) {
            $end = $_POST['end'];
        }

        $script = "python3 /var/www/html/plugins/Autotagging/libraries/autotagging/auto_tag.py";
        $args = "http://allenarchive-dev.iac.gatech.edu/ 030c516f3f818bb10793ff6c965489c69647129d -s {$start} -e {$end}";
        $command = "$script $args";
        $out = [];
        $status = [];
        exec($command, $out, $status);
        foreach ($out as $value) {
            echo ("$value\n");
        }
    }
}
