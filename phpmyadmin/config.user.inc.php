<?php
/**
 * phpMyAdmin sample configuration, you can use it as base for
 * manual configuration. For easier setup you can use setup/
 *
 * All directives are explained in documentation in the doc/ folder
 * or at <https://docs.phpmyadmin.net/>.
 */

declare(strict_types=1);

$i=0;
$i++;
$cfg['Servers'][$i]['auth_type']     = 'signon';
$cfg['Servers'][$i]['SignonSession'] = 'SignonSession';
// $cfg['Servers'][$i]['SignonURL']     = 'examples/back-to-ecp.php';
// $cfg['Servers'][$i]['SignonURL']     = 'examples/custom-signon-url.php';

// $cfg['Servers'][$i]['SignonScript'] = 'examples/custom-signon.php';

// $cfg['Servers'][$i]['LogoutURL'] = 'http://localhost:2325/';
$cfg['Servers'][$i]['LogoutURL'] = '/ehm-authentication/logout.php';