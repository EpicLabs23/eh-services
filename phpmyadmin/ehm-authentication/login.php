<?php

declare(strict_types=1);

$curl = curl_init();

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
    setcookie('logout_redirect_url', $_GET['logout_redirect_url']);
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

    // var_dump($apiResponse);
    // die;

    curl_close($curl);

    $apiData = json_decode($apiResponse, true);
    // var_dump($apiData);die;

    $dbUser = $apiData['dbUser'];
    $dbPass = $apiData['dbPass'];
    $dbHost = $apiData['dbHost'];
    $dbPort = $apiData['dbPort'];
    $host_verbose = $apiData['host_verbose'];

    /* Store there credentials */
    $_SESSION['PMA_single_signon_user'] = $dbUser;
    $_SESSION['PMA_single_signon_password'] = $dbPass;
    $_SESSION['PMA_single_signon_host'] = $dbHost;
    $_SESSION['PMA_single_signon_port'] = $dbPort;
    /* Update another field of server configuration */
    $_SESSION['PMA_single_signon_cfgupdate'] = ['verbose' => $host_verbose];
    $_SESSION['PMA_single_signon_HMAC_secret'] = hash('sha1', uniqid(strval(random_int(0, mt_getrandmax())), true));
    $id = session_id();
    /* Close that session */
    @session_write_close();
    /* Redirect to phpMyAdmin (should use absolute URL here!) */
    header('Location: ../index.php');
}
