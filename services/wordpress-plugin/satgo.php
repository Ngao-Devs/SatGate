<?php
/**
 * Plugin Name: SatGo - Lightning Pay-to-Submit
 * Description: Simple integration for SatGo Lightning-gated forms.
 * Version: 1.0.0
 * Author: SatGo
 */

if (!defined('ABSPATH')) {
    exit;
}

add_action('admin_menu', 'satgo_menu');
function satgo_menu()
{
    add_options_page('SatGo Settings', 'SatGo', 'manage_options', 'satgo', 'satgo_settings_page');
}

function satgo_settings_page()
{
    ?>
    <div class="wrap">
        <h1>SatGo Settings</h1>
        <form method="post" action="options.php">
            <?php
            settings_fields('satgo_settings');
    do_settings_sections('satgo');
    submit_button();
    ?>
        </form>
    </div>
    <?php
}

add_action('admin_init', 'satgo_settings_init');
function satgo_settings_init()
{
    register_setting('satgo_settings', 'satgo_site_id');
    register_setting('satgo_settings', 'satgo_amount');

    add_settings_section('satgo_main', 'Main Settings', null, 'satgo');

    add_settings_field('satgo_site_id', 'Site ID', function () {
        $val = get_option('satgo_site_id');
        echo '<input type="text" name="satgo_site_id" value="' . esc_attr($val) . '" class="regular-text">';
    }, 'satgo', 'satgo_main');

    add_settings_field('satgo_amount', 'Default Amount (SATS)', function () {
        $val = get_option('satgo_amount', '100');
        echo '<input type="number" name="satgo_amount" value="' . esc_attr($val) . '">';
    }, 'satgo', 'satgo_main');
}

add_action('wp_footer', 'satgo_inject_script');
function satgo_inject_script()
{
    $site_id = get_option('satgo_site_id');
    $amount = get_option('satgo_amount', '100');
    if (!$site_id) {
        return;
    }

    // Replace with your production URL
    echo '<script src="http://localhost:3000/widget.js" data-site-id="' . esc_attr($site_id) . '" data-amount="' . esc_attr($amount) . '"></script>';
}
