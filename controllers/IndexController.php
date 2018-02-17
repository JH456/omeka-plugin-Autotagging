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
    {
        // If we want, we can assign variables here. Like, if we assign nodes
        // to something, in views/scripts/browse.php we will have $nodes
        // defined with that value.
        // If we want to do any server side processing (like database queries)
        // we should do them in this file, and then pass the results to the
        // view like this:
        // You can take a look at
        // application/libraries/Omeka/Controller/AbstractActionController.php
        // to see how to run database queries and the like, but I think we
        // are mostly going to be doing stuff with elastic search here
    }
}
