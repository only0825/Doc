//
//  testc.c
//  mymacos
//
//  Created by WH37 on 2022/9/15.
//

#include "testc.h"

void haha() {

    int ret = 0;
    char errors[1024];
    
    AVFormatContext *fmt_ctx = NULL;
    AVDictionary *options = NULL;
    
    char *devicename = ":0"; // [[vide device]:[audio device]]
    
    av_log_set_level(AV_LOG_DEBUG);
    
    // register audio device
    avdevice_register_all();
    
    AVInputFormat *iformat = av_find_input_format("avfoundation"); // get format
    
    ret = avformat_open_input(&fmt_ctx, devicename, iformat, &options);
    if (ret < 0) {
        av_strerror(ret, errors, 1024);
        fprintf(stderr, "Failed to open audio device, [%d] %s\n", ret, errors);
        return;
    }

    av_log(NULL, AV_LOG_DEBUG, "hello, world!\n");
    
    return;
}
