#ifndef __prism_h__
#define __prism_h__

#include <opencv/cv.h>
#include <opencv/highgui.h>
#include <opencv2/imgproc/types_c.h>
#include <stdio.h>
#include <turbojpeg.h>

typedef struct {
  unsigned char* buffer;
  unsigned long size;
  CvMat* _mat;
} PrismEncoded;

PrismEncoded* prismEncodeJPEG(IplImage* img, int quality);
PrismEncoded* prismEncodePNG(IplImage* img, int compression);

void prismRelease(PrismEncoded* enc);

IplImage* prismDecode(void* data, unsigned int dataSize);

#endif
