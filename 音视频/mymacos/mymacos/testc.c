//
//  testc.c
//  mymacos
//
//  Created by WH37 on 2022/9/15.
//

#include "testc.h"
static int rec_status = 0;

void set_status(int status) {
    rec_status = status;
}

void rec_audio() {

    int ret = 0;
    char errors[1024] = {0, };
    
    // ctx
    AVFormatContext *fmt_ctx = NULL;
    AVDictionary *options = NULL;
    
    // packet
    int count = 0;
    AVPacket pkt;
    
    // [[vide device]:[audio device]]
    char *devicename = ":0";
    
    // set log level
    av_log_set_level(AV_LOG_DEBUG);
    
    // start record
    rec_status = 1;
    
    // register audio device
    avdevice_register_all();
    
    // get format
    AVInputFormat *iformat = av_find_input_format("avfoundation");
    
    // open device
    if((ret = avformat_open_input(&fmt_ctx, devicename, iformat, &options)) < 0 ){
        av_strerror(ret, errors, 1024);
        fprintf(stderr, "Failed to open video device, [%d]%s\n", ret, errors);
        return;
    }
    
    // create file
    char *out = "/Users/wh37/Documents/Doc/音视频/audio.pcm";
    FILE *outfile = fopen(out, "wb+"); // w写 b二进制 +不存在则创建

    //av_init_packet(&pkt); // 现在不需要初始化了
    
    // read data from device
    ret = 0;
    while(ret == 0 && rec_status) {
        ret = av_read_frame(fmt_ctx, &pkt);
//        av_log(NULL, AV_LOG_INFO, "pkt size is %d(%p), count=%d \n", pkt.size, pkt.data, count);
        printf("ret = %d at 第%d次\n", ret, count);
        // - 35 应该是设备暂时还没准备好，休息一秒继续尝试
        if (ret == -35) {
            sleep(1);
            ret = 0;
            // 释放音视频包
            av_packet_unref(&pkt);
            continue;
        }
        
        if (ret < 0) {
            // 释放音视频包
            av_packet_unref(&pkt);
            break;
        }
        
        // write file
        fwrite(pkt.data, pkt.size, 1, outfile);
        fflush(outfile);
        
        printf("pkt size is %d at 第%d次 \n", pkt.size, count);
        count ++;
        // 释放音视频包
        av_packet_unref(&pkt);
    }
    
    // close file
    fclose(outfile);
    
    // close device and release ctx
    avformat_close_input(&fmt_ctx);
    
    av_log(NULL, AV_LOG_DEBUG, "finish\n");
    
    return;
}
