<?php

// echo "Cookie: <pre>";
// print_r($_COOKIE);
$logout_url = $_COOKIE['logout_redirect_url'];
setcookie('logout_redirect_url', "", time() - 3600);
// echo $logout_url;
header("Location: $logout_url");