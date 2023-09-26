<?php

/**
 * Single signon for phpMyAdmin
 *
 * This is just example how to use session based single signon with
 * phpMyAdmin, it is not intended to be perfect code and look, only
 * shows how you can integrate this functionality in your application.
 */

declare(strict_types=1);

$curl = curl_init();

// curl_setopt_array($curl, array(
//   CURLOPT_URL => 'http://host.docker.internal:2326/account/mysqlcred',
//   CURLOPT_RETURNTRANSFER => true,
//   CURLOPT_ENCODING => '',
//   CURLOPT_MAXREDIRS => 10,
//   CURLOPT_TIMEOUT => 0,
//   CURLOPT_FOLLOWLOCATION => true,
//   CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
//   CURLOPT_CUSTOMREQUEST => 'GET',
//   CURLOPT_HTTPHEADER => array(
//     'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InJ5bWFtIiwic3ViIjoyOCwiaWF0IjoxNjkzMTg1MDc1LCJleHAiOjE2OTMyNzE0NzV9.ovKvKwfVMqJbIyZk3lu_CtAMx0gNdBDA5Wh2KZ77wK0',
//     'Cookie: pma_lang=en'
//   ),
// ));

// $response = curl_exec($curl);

// curl_close($curl);
// echo $response;
// die;

if (isset($_GET['access_token'])) {
    /* Use cookies for session */
    ini_set('session.use_cookies', 'true');
    /* Change this to true if using phpMyAdmin over https */
    $secureCookie = false;
    /* Need to have cookie visible from parent directory */
    session_set_cookie_params(0, '/', '', $secureCookie, true);
    /* Create signon session */
    $sessionName = 'SignonSession';
    session_name($sessionName);
    // Uncomment and change the following line to match your $cfg['SessionSavePath']
    //session_save_path('/foobar');
    @session_start();

    $scope = isset($_GET['scope']) ? $_GET['scope'] : null;
    $token = $_GET['access_token'];
    $ehmurl = 'http://host.docker.internal:2326';

    if ($scope == 'ehm') {
        $apiUrl = "$ehmurl/user/mysqlcred";
    } else {
        $apiUrl = "$ehmurl/account/mysqlcred";
    }

    $curl = curl_init();

    curl_setopt_array($curl, array(
        CURLOPT_URL => $apiUrl,
        CURLOPT_RETURNTRANSFER => true,
        CURLOPT_ENCODING => '',
        CURLOPT_MAXREDIRS => 10,
        CURLOPT_TIMEOUT => 0,
        CURLOPT_FOLLOWLOCATION => true,
        CURLOPT_HTTP_VERSION => CURL_HTTP_VERSION_1_1,
        CURLOPT_CUSTOMREQUEST => 'GET',
        CURLOPT_HTTPHEADER => array(
            "Authorization: Bearer $token",
            'Cookie: pma_lang=en'
        ),
    ));

    $apiResponse = curl_exec($curl);

    curl_close($curl);

    $apiData = json_decode($apiResponse, true);
// var_dump($apiData);die;

    $dbUser = $apiData['dbUser'];
    $dbPass = $apiData['dbPass'];
    $dbHost = $apiData['dbHost'];
    $dbPort = $apiData['dbPort'];

    /* Store there credentials */
    $_SESSION['PMA_single_signon_user'] = $dbUser;
    $_SESSION['PMA_single_signon_password'] = $dbPass;
    $_SESSION['PMA_single_signon_host'] = $dbHost;
    $_SESSION['PMA_single_signon_port'] = $dbPort;
    /* Update another field of server configuration */
    $_SESSION['PMA_single_signon_cfgupdate'] = ['verbose' => 'Signon test'];
    $_SESSION['PMA_single_signon_HMAC_secret'] = hash('sha1', uniqid(strval(random_int(0, mt_getrandmax())), true));
    $id = session_id();
    /* Close that session */
    @session_write_close();
    /* Redirect to phpMyAdmin (should use absolute URL here!) */
    header('Location: ../index.php');
}

/**
 * This function returns username and password.
 *
 * It can optionally use configured username as parameter.
 *
 * @param string $user User name
 *
 * @return array<int,string>
 */
function get_login_credentials(string $user): array
{
    if (isset($_SESSION['PMA_single_signon_user'])) {
        return [$_SESSION['PMA_single_signon_user'], $_SESSION['PMA_single_signon_password']];
    } else {
        echo ("All empty");
        die;
    }
}
