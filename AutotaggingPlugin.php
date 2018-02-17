<?php
/**
 * Auto Tagging
 *
 * @license http://www.gnu.org/licenses/gpl-3.0.txt GNU GPLv3
 */

/**
 * Auto Tagging plugin.
 */
class AutotaggingPlugin extends Omeka_Plugin_AbstractPlugin
{
    /**
     * @var array Hooks for the plugin.
     */
    protected $_hooks = array(
        'config_form',
        'public_items_show',
        'admin_items_show',
        'public_items_browse_each',
        'admin_items_browse_simple_each'
    );

    /**
     * @var array Filters for the plugin.
     */
    protected $_filters = array('admin_navigation_main', 'public_navigation_main');

    /**
     * Display the plugin config form.
     */
    public function hookConfigForm()
    {
        require dirname(__FILE__) . '/config_form.php';
    }

    public function hookPublicItemsShow()
    {
        echo '<strong>public items show</strong>';
    }
    
    public function hookAdminItemsShow()
    {
        echo '<strong>admin items show</strong>';
    }
    
    public function hookPublicItemsBrowseEach()
    {
        echo '<strong>public items browse each</strong>';
    }
    
    public function hookAdminItemsBrowseSimpleEach()
    {
        echo '<strong>admin items browse simple each</strong>';
    }

    /**
     * Add the Auto Tagging link to the admin main navigation.
     * 
     * @param array Navigation array.
     * @return array Filtered navigation array.
     */
    public function filterAdminNavigationMain($nav)
    {
        $nav[] = array(
            'label' => __('Auto Tagging'),
            'uri' => url('autotagging'),
        );
        return $nav;
    }

    /**
     * Add the pages to the public main navigation options.
     * 
     * @param array Navigation array.
     * @return array Filtered navigation array.
     */
    public function filterPublicNavigationMain($nav)
    {
        $nav[] = array(
            'label' => __('Auto Tagging'),
            'uri' => url('autotagging'),
        );
        return $nav;
    }
}
