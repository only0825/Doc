//
//  testc.h
//  mymacos
//
//  Created by WH37 on 2022/9/15.
//

#ifndef testc_h
#define testc_h

#include <stdio.h>
#include "libavutil/avutil.h"
#include "libavdevice/avdevice.h"
#include "libavformat/avformat.h"
#include <unistd.h>

void rec_audio(void);
void set_status(int status);

#endif /* testc_h */
