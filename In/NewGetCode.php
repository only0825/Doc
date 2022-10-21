<?php

$servername = "101.32.18.63";
$username = "lengjing";
$password = "sw8ct7hkL3rwRwtz";
$dbname = "lengjing";
$conn = mysqli_connect($servername, $username, $password, $dbname);
if (!$conn)
  {
  die('Could not connect: ' . mysqli_connect_error());
  }
$date = 666;
$num = mt_rand(100000000000, 9999999999999);
$sql = "INSERT INTO ipaccode (accode) VALUES ('". $num ."')";

if ($conn->query($sql) === TRUE) {
    echo "新记录插入成功";
} else {
    echo "Error: " . $sql . "<br>" . $conn->error;
}
?>