<?php
$servername = "101.32.18.63";
$username = "lengjing";
$password = "sw8ct7hkL3rwRwtz";
$con = mysqli_connect($servername, $username, $password);
if (!$con)
  {
  die('Could not connect: ' . mysqli_connect_error());
  }
$date = 666;
$num = mt_rand(1000000000000, 9999999999999);
mysqli_query($con,"INSERT INTO ipaccode (accode) 
VALUES ('$num')");
echo "$num";
?>