//
//  testc.c
//  mymacos
//
//  Created by WH37 on 2022/9/15.
//

#include "testc.h"
#include <string.h>

static int rec_status = 0;

void set_status(int status) {
    rec_status = status;
}

SwrContext* init_swr(void){
    
    SwrContext *swr_ctx = NULL;
    
    //channel, number/
    swr_ctx = swr_alloc_set_opts(NULL,                //ctx
                                 AV_CH_LAYOUT_STEREO, //输出channel布局
                                 AV_SAMPLE_FMT_S16,   //输出的采样格式
                                 44100,               //采样率
                                 AV_CH_LAYOUT_STEREO, //输入channel布局
                                 AV_SAMPLE_FMT_FLT,   //输入的采样格式
                                 44100,               //输入的采样率
                                 0, NULL);
    
    if(!swr_ctx){
        
    }
    
    if(swr_init(swr_ctx) < 0){
        
    }
    
    return swr_ctx;
}

void rec_audio() {
    
    int ret = 0;
    char errors[1024] = {0, };
    
    //重采样缓冲区
    uint8_t **src_data = NULL;
    int src_linesize = 0;
    
    uint8_t **dst_data = NULL;
    int dst_linesize = 0;
    
    // ctx
    AVFormatContext *fmt_ctx = NULL;
    AVDictionary *options = NULL;
    
    // packet
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
    const AVInputFormat *iformat = av_find_input_format("avfoundation");
    
    // open device
    if((ret = avformat_open_input(&fmt_ctx, devicename, iformat, &options)) < 0 ){
        av_strerror(ret, errors, 1024);
        fprintf(stderr, "Failed to open video device, [%d]%s\n", ret, errors);
        return;
    }
    
    // create file
    char *out = "/Users/wh37/workspace/sound/audio2.pcm";
    FILE *outfile = fopen(out, "wb+"); // w写 b二进制 +不存在则创建
    
    // 打开编码器
    // avcodec_find_encoder(AV_CODEC_ID_AAC)
    const AVCodec *codec = avcodec_find_encoder_by_name("libfdk_aac");
    
    // 创建 codec 上下文
    AVCodecContext *codec_ctx = avcodec_alloc_context3(codec);
    
    codec_ctx->sample_fmt = AV_SAMPLE_FMT_S16; // 输入音频的采样大小
    codec_ctx->channel_layout = AV_CH_LAYOUT_STEREO; // 输入音频的 channel layout
    
    SwrContext* swr_ctx = init_swr();
    
    // 4096/4 = 1024/2 = 512
    // 创建输入缓冲区
    av_samples_alloc_array_and_samples(&src_data,           // 输出缓冲区地址
                                       &src_linesize,       // 缓冲区的大小
                                       2,                   // 通道个数
                                       512,                 // 单通道采样个数
                                       AV_SAMPLE_FMT_FLT,   // 采样格式
                                       0);
    
    // 创建输出缓冲区
    av_samples_alloc_array_and_samples(&dst_data,           // 输出缓冲区地址
                                       &dst_linesize,       // 缓冲区的大小
                                       2,                   // 通道个数
                                       512,                 // 单通道采样个数
                                       AV_SAMPLE_FMT_S16,   // 采样格式
                                       0);
    
    
    //av_init_packet(&pkt); // 现在不需要初始化了
    
    // read data from device
    ret = 0;
    int count = 0;
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
        
        // 进行内存拷贝，按字节拷贝的
        memcpy((void*)src_data[0], (void*)pkt.data, pkt.size);
        
        //重采样
        swr_convert(swr_ctx,                    //重采样的上下文
                    dst_data,                   //输出结果缓冲区
                    512,                        //每个通道的采样数
                    (const uint8_t **)src_data, //输入缓冲区
                    512);                       //输入单个通道的采样数
        
        // write file
        //fwrite(pkt.data, pkt.size, 1, outfile);
        fwrite(dst_data[0], 1, 512, outfile);
        fflush(outfile);
        
        printf("pkt size is %d at 第%d次 \n", pkt.size, count);
        count ++;
        // 释放音视频包
        av_packet_unref(&pkt);
    }
    
    // close file
    fclose(outfile);
    
    // 释放输入输出缓冲区
    if (src_data) {
        av_freep(&src_data[0]);
    }
    av_freep(src_data);
    
    if (dst_data) {
        av_freep(&dst_data[0]);
    }
    av_freep(dst_data);
    
    // 释放重采样的上下文
    swr_free(&swr_ctx);
    
    // close device and release ctx
    avformat_close_input(&fmt_ctx);
    
    av_log(NULL, AV_LOG_DEBUG, "finish\n");
    
    return;
}
