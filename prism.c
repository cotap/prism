#include <prism.h>

int isRecoverableError(char* errorStr) {
  if (!strcmp(errorStr, "Invalid SOS parameters for sequential JPEG") ||
      !strcmp(errorStr, "Premature end of JPEG file")) {
    return 1;
  }
  return 0;
}

void printError(char* errorStr) {
  fprintf(stderr, "prism: %s\n", errorStr);
  fflush(stderr);
}

IplImage* prismDecode(void* data, unsigned int dataSize) {
  int err;
  IplImage* iplImage;
  tjhandle jpeg = tjInitDecompress();
  int width, height, subsamp, colorspace;

  // attempt to decode JPEG header
  err = tjDecompressHeader3(jpeg, (unsigned char*)data, dataSize, &width, &height, &subsamp, &colorspace);
  if (err) {
    tjDestroy(jpeg);

    // fall back to OpenCV decoding
    CvMat* cvMat = cvCreateMatHeader(1, dataSize, CV_8UC1);
    cvSetData(cvMat, data, dataSize);
    iplImage = cvDecodeImage(cvMat, CV_LOAD_IMAGE_UNCHANGED);
    cvReleaseMat(&cvMat);
    if (iplImage) {
      return iplImage;
    }

    printError(tjGetErrorStr());
    return NULL;
  }

  int channels = 3, pixelFmt = TJPF_BGR;
  if (colorspace == TJCS_GRAY) {
    channels = 1;
    pixelFmt = TJPF_GRAY;
  }

  unsigned char* buffer = cvAlloc(width * height * channels);
  err = tjDecompress2(
          jpeg, (unsigned char*)data, dataSize, buffer, 0, 0, 0, pixelFmt, TJFLAG_FASTDCT
        );
  tjDestroy(jpeg);

  if (err) {
    char* errorStr = tjGetErrorStr();
    if (!isRecoverableError(errorStr)) {
      cvFree(&buffer);
      printError(errorStr);
      return NULL;
    }
  }

  iplImage = cvCreateImageHeader(cvSize(width, height), IPL_DEPTH_8U, channels);
  cvSetData(iplImage, buffer, width * channels);
  return iplImage;
}

PrismEncoded* prismEncodeJPEG(IplImage* img, int quality) {
  int pixFmt, subsamp;

  switch (img->nChannels) {
  case 1:
    pixFmt = TJPF_GRAY;
    subsamp = TJSAMP_GRAY;
    break;
  case 4:
    pixFmt = TJPF_BGRA;
    subsamp = TJSAMP_420;
    break;
  default:
    pixFmt = TJPF_BGR;
    subsamp = TJSAMP_420;
  }

  int err;
  CvSize size = cvGetSize(img);
  tjhandle jpeg = tjInitCompress();
  PrismEncoded* enc = calloc(1, sizeof(PrismEncoded));

  err = tjCompress2(
          jpeg, (unsigned char*)img->imageData, size.width, img->widthStep,
          size.height, pixFmt, &enc->buffer, &enc->size, subsamp, quality, 0
        );
  tjDestroy(jpeg);

  if (err) {
    printError(tjGetErrorStr());
    prismRelease(enc);
    return NULL;
  }

  return enc;
}

PrismEncoded* prismEncodePNG(IplImage* img, int compression) {
  int encodeParams[5] = {
    CV_IMWRITE_PNG_COMPRESSION, compression,
    CV_IMWRITE_PNG_STRATEGY, 0,
    0,
  };

  CvMat* cvMat = cvEncodeImage(".png", img, encodeParams);
  if (!cvMat) {
    return NULL;
  }

  PrismEncoded* enc = malloc(sizeof(PrismEncoded));
  enc->buffer = cvMat->data.ptr;
  enc->size = cvMat->cols;
  enc->_mat = cvMat;

  return enc;
}

void prismRelease(PrismEncoded* enc) {
  if (enc->_mat) {
    cvReleaseMat(&enc->_mat);
  } else if (enc->buffer) {
    tjFree(enc->buffer);
  }
  free(enc);
}
