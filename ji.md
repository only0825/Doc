![telegram-cloud-photo-size-5-6249201528082641150-x](/Users/wh37/Library/Group Containers/6N38VWS5BX.ru.keepcoder.Telegram/stable/account-603578246170668601/postbox/media/telegram-cloud-photo-size-5-6249201528082641150-x.jpg)



ufcnation.net

有个趋势图，我想要做成这样的，但是不知道数据怎么处理，你打开他的控制台看到trends这个请求，就是数据结构



side: -1 CANCELLED 灰色

side: 1 MERON 红色

side: 2 WALA 蓝色

side: 3 DRAW 绿色

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

