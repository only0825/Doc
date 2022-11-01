![telegram-cloud-photo-size-5-6249201528082641150-x](/Users/wh37/Library/Group Containers/6N38VWS5BX.ru.keepcoder.Telegram/stable/account-603578246170668601/postbox/media/telegram-cloud-photo-size-5-6249201528082641150-x.jpg)



ufcnation.net

有个趋势图，我想要做成这样的，但是不知道数据怎么处理，你打开他的控制台看到trends这个请求，就是数据结构



side: -1 CANCELLED 灰色

side: 1 MERON 红色

side: 2 WALA 蓝色

side: 3 DRAW 绿色



表字段 game_num  winner(1红色 2蓝色 3平局 4取消)

```php
    // 有新数据时调用
    public function updateTrends()
    {
        // 这里是模拟数据，实际是要传参
        $winner = 'WALA';

        // 如果没有winner返回错误提示
        if (empty($winner)) {
            die;
        }

        $data = Cache::store('redis')->get('trends');

        if (!empty($data)) {
            $num = Cache::store('redis')->get('trendsNum');
            $data = json_decode($data, true);
            for ($i = 0; $i < 6; $i++) {
                $len = count($data[$i]);
                for ($j = 0; $j < $len; $j++) {
                    // 当前值不为null且下一列为空
                    if ($data[$i][$j]['num'] != NULL) {
                        if (empty($data[$i][$j + 1]['num'])) {
                            $side = $this->getSide($winner);
                            // 如果side相同就增加行 否则就换列
                            if ($side == $data[$i][$j]['side']) {
                                // 如果6行都有数据 则放到下一列
                                if ($i == 5) {
                                    $num = $num + 1;
                                    $data[0][$j + 1]['num'] = $num;
                                    $data[0][$j + 1]['side'] = $side;
                                } else {
                                    if (empty($data[$i + 1][$j])) {
                                        $data[$i + 1][$j]['num'] = null;
                                        $data[$i + 1][$j]['side'] = 0;
                                    }
                                    if ($data[$i + 1][$j]['num'] == NULL) {
                                        $num = $num + 1;
                                        $data[$i + 1][$j]['num'] = $num;
                                        $data[$i + 1][$j]['side'] = $side;
                                    } else {
                                        continue;
                                    }
                                }
                            } else {
                                $num = $num + 1;
                                $data[0][$j + 1]['num'] = $num;
                                $data[0][$j + 1]['side'] = $side;
                            }
                            goto end;
                        }
                    }
                }
            }
        } else {
            // 初始化操作
            $data = [];
            $num = 1;
            for ($i = 0; $i < 6; $i++) {
                $data[$i] = [];
                for ($j = 0; $j < 10; $j++) {
                    if ($i == 0 && $j == 0) {
                        $data[$i][$j]['num'] = $num;
                        $data[$i][$j]['side'] = $this->getSide($winner);
                    } else {
                        $data[$i][$j]['num'] = null;
                        $data[$i][$j]['side'] = 0;
                    }
                }
            }
        }

        end:
        Cache::store('redis')->set('trends', json_encode($data));
        Cache::store('redis')->set('trendsNum', $num);

        return $this->ReturnSuccess($data);
    }

    public function getSide($winner)
    {
        $side = 0;
        switch ($winner) {
            case "MERON":
                $side = 1;
                break;
            case "WALA":
                $side = 2;
                break;
            case "DRAW":
                $side = 3;
                break;
            case "CANCELLED":
                $side = -1;
                break;
        }
        return $side;
    }
```



```php
   public function getTrends()
    {
        // 查询最新的N条数据
        $res = Db::table('ji_test')->field(["game_num", "winner"])->order('created_at', 'ASC')->limit(100)->select()->toArray();

        // 初始化数组
        $data = [];
        for ($i = 0; $i < 6; $i++) {
            for ($j = 0; $j < count($res); $j++) {
                $data[$i][$j]['num'] = null;
                $data[$i][$j]['winner'] = 0;
            }
        }

        // 将$res的数据存入$data
        foreach ($res as $arr) {
            $winner = $arr['winner'];
            $num = $arr['game_num'];
            $data = $this->saveToData($data, $winner, $num);
        }

        // 删除空数据
        $len = 0;
        for ($i = 0; $i < count($data[0]); $i++) {
            if ($data[0][$i]['num'] != null) {
                $len++;
            }
        }
        for ($i = 0; $i < count($data); $i++) {
            $data[$i] = array_slice($data[$i], 0, $len);
        }

        echo json_encode($data);
    }

    public function saveToData($data, $winner, $num)
    {
        if ($data[0][0]['num'] == null) {
            $data[0][0]['num'] = $num;
            $data[0][0]['winner'] = $winner;
            return $data;
        }

        for ($i = 0; $i < 6; $i++) {
            for ($j = 0; $j < count($data[$i]); $j++) {
                // 当前行列不为null，并且为最后一列
                if ($data[$i][$j]['num'] != null && empty($data[$i][$j + 1]['num'])) {
                    // 如果已经到最后一行，则在新的一列存数据
                    if ($i == 5) {
                        $data[0][$j + 1]['num'] = $num;
                        $data[0][$j + 1]['winner'] = $winner;
                        return $data;
                    }
                    // 如果下一行的数据是null
                    if ($data[$i + 1][$j]['num'] == null) {
                        if ($data[$i][$j]['winner'] == $winner) {
                            // 1.当前行列的winner等于新的winner 就在下一行存数据
                            $data[$i + 1][$j]['num'] = $num;
                            $data[$i + 1][$j]['winner'] = $winner;
                        } else {
                            // 2.当前行列的winner 不等于新的winner 就在下一列的第一行存数据
                            $data[0][$j + 1]['num'] = $num;
                            $data[0][$j + 1]['winner'] = $winner;
                        }
                        return $data;
                    }
                }
            }
        }
        return $data;
    }

```



```

//        for ($i = 1; $i < 200; $i++) {
//            $where['game_num'] = $i;
//            $arr = [1,1,1,1,1,1,1,1,2,2,2,2,2,2,2,2,3,4];
//            $index = array_rand($arr);
////            $data['created_at'] = date('Y-m-d H:i:s');
//            $data['winner'] = $arr[$index];
//            Db::table('ji_test')->where($where)->update($data);
////            sleep(1);
//        }
```

