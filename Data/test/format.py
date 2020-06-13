# Usage: python format.py labels.csv
import os
import numpy as np
import pandas as pd
import sys
import cv2
import copy
import imutils
import math

class_name = "tag"
class_id = 0

# Create dataframes for csvs
annot = pd.read_csv(sys.argv[1])

# Dropping unneeded columns from the bounding box files
cw = os.getcwd()
annot["filename"] = [cw + "/" + fn for fn in annot["filename"]]
names = annot["filename"].tolist()

# names.to_csv("lists.txt", index=False)

annot = annot.drop(["filename", "class"], axis=1)

# Calculate x, y, width, height relative to image height

x = [coord / dim for coord, dim in zip(annot["xmin"], annot["width"])]
y = [coord / dim for coord, dim in zip(annot["ymin"], annot["height"])]

width = [(xmax - xmin) / dim for xmax, xmin, dim in zip(annot["xmax"], annot["xmin"], annot["width"])]
height = [(ymax - ymin) / dim for ymax, ymin, dim in zip(annot["ymax"], annot["ymin"], annot["height"])]
classes = np.full((annot.shape[0]), 0)

annot = annot.drop(["width", "height", "xmin", "ymin", "xmax", "ymax"], axis=1)

annot["object-class"] = classes
annot["x"] = x
annot["y"] = y
annot["width"] = width
annot["height"] = height

master = pd.DataFrame.from_records([], columns=["filename", "width", "height",
                                            "class", "xmin", "ymin", "xmax", "ymax"])


for filename, entry in zip(names, annot.itertuples()):
    data = [entry[1:]]
    df = pd.DataFrame.from_records(data, columns=["object-class", "x", "y",
        "width", "height"])
    df.to_csv(filename[:len(filename) - 4] + ".txt", index=False, header=False,
            sep=" ")
    # Augmenting
    img = cv2.imread(filename)

    def identity(data, level):
        return data

    def apply_transform(transform, name, sizes, datatrans = identity):
        for level in sizes:
            newname = filename[:len(filename) - 4] + name + str(level) + filename[len(filename) - 4:]
            transformed = transform(img, level)
            names.append(newname)
            cv2.imwrite(newname, transformed)
            transformeddata = datatrans(data, level)

            transformname = filename[:len(filename) - 4] + name + str(level) + ".JPG"
            toadd = pd.DataFrame([[transformname, len(transformed[0]),
                                   len(transformed), "ARTag",
                                   transformeddata[0][1] * len(transformed[0]),
                                   transformeddata[0][2] * len(transformed),
                                   (transformeddata[0][1] + transformeddata[0][3]) * len(transformed[0]),
                                   (transformeddata[0][2] + transformeddata[0][4]) * len(transformed)
                                   ]], columns=["filename", "width", "height",
                                                        "class", "xmin", "ymin", "xmax", "ymax"])
            global master
            master = master.append(toadd, ignore_index = True)

    # vflip
    def vflipdata(data, level):
        vflipdata = copy.deepcopy(data)
        vflipdata[0] = (vflipdata[0][0], 1 - vflipdata[0][1] - vflipdata[0][3],
                        vflipdata[0][2], vflipdata[0][3], vflipdata[0][4])
        return data
    apply_transform(lambda img, level : cv2.flip(img, 0), "vflip", [0], vflipdata)

    #bright
    apply_transform(lambda img, level : cv2.add(img, np.array([level])),
                    "bright", [-40.0, -20.0, 20.0, 40.0])

    # hFlip
    def hflipdata(data, level):
        hflipdata = copy.deepcopy(data)
        hflipdata[0] = (hflipdata[0][0], 1 - hflipdata[0][1] - hflipdata[0][3],
                        hflipdata[0][2], hflipdata[0][3], hflipdata[0][4])
        return data
    apply_transform(lambda img, level : cv2.flip(img, 1),
                    "hflip", [0],  hflipdata)

    # saturation
    def satimg(img, level):
        hsvImg = cv2.cvtColor(img,cv2.COLOR_BGR2HSV)
        hsvImg[...,1] = hsvImg[...,1]*level
        satimg=cv2.cvtColor(hsvImg,cv2.COLOR_HSV2BGR)
        return satimg
    apply_transform(satimg, "sat", [.5, .75, 1.25, 1.5])

    # salt & pepper
    def salt_pepper(img, level):
        rand = cv2.randn(np.zeros_like(img), (0, 0, 0), (level, level, level))
        noise = cv2.add(img, rand)
        return noise
    apply_transform(salt_pepper, "noise", [40.0, 80.0, 120.0])

    # blur
    apply_transform(lambda img, ksize: cv2.GaussianBlur(img, (ksize, ksize), 0),
                    "blur", [5, 9, 13])

    # motion blur horizontal
    def hblur(img, ksize):
        kernel_h = np.zeros((ksize, ksize))
        kernel_h[int((ksize - 1)/2), :] = np.ones(ksize)
        kernel_h /= ksize
        hblurred = cv2.filter2D(img, -1, kernel_h)
        return hblurred
    apply_transform(hblur, "hblur", [3, 7, 21],)

    # motion blur vertical
    def vblur(img, ksize):
        kernel_v = np.zeros((ksize, ksize))
        kernel_v[int((ksize - 1)/2), :] = np.ones(ksize)
        kernel_v /= ksize
        vblurred = cv2.filter2D(img, -1, kernel_v)
        return vblurred
    apply_transform(vblur, "vblur", [3, 7, 21])

    # zoom
    def zoom(img, zoomlevel):
        minx = int(max(data[0][1] - data[0][3] / zoomlevel, 0) * len(img[0]))
        miny = int(max(data[0][2] - data[0][4] / zoomlevel, 0) * len(img))

        maxx = int(min(data[0][1] + data[0][3] * (1 + 1 / zoomlevel), 1) *
                   len(img[0]))
        maxy = int(min(data[0][2] + data[0][4] * (1 + 1 / zoomlevel), 1) *
                   len(img))
        cropimg = img[miny:maxy, minx:maxx]
        return cropimg

    def zoomdata(data, level):
        bufpercent = (float(1 / level) / (1 + 2 / level))
        zoomdata = copy.deepcopy(data)
        zoomdata[0] = (zoomdata[0][0], bufpercent,
                        bufpercent, 1 - 2 * bufpercent, 1 - 2 * bufpercent)
        return zoomdata
    apply_transform(zoom, "zoom", [2, 5, 15],
                    zoomdata)

    # rotation
    def rotatedata(data, level):
        midx = (data[0][1] + data[0][3]/2 - .5) * len(img[0])
        midy = (data[0][2] + data[0][4]/2 - .5) * len(img)

        cos = math.cos(level / 180 * math.pi)
        sin = math.sin(level / 180 * math.pi)
        newmidx = (cos * midx - sin * midy) / (abs(sin) + abs(cos)) + len(img[0]) / 2
        newmidy = (sin * midx + cos * midy) / (abs(sin) + abs(cos)) + len(img) / 2

        xscale = (len(img[0]))
        yscale = (len(img))

        minx = newmidx - (data[0][3]/2) * len(img[0]) / (abs(sin) + abs(cos))
        miny = newmidy - (data[0][4]/2) * len(img) / (abs(sin) + abs(cos))
        width = (newmidx - minx) * 2
        height = (newmidy - miny) * 2

        rotatedata = copy.deepcopy(data)
        rotatedata[0] = (rotatedata[0][0], minx / xscale, miny / yscale, width/xscale, height/yscale)
        return rotatedata

    apply_transform(lambda image, level: imutils.rotate_bound(image, level),
                    "rotate", range(0, 360, 30),  rotatedata)
    # zoom out
    def zoomout(img, zoomlevel):
        dst = np.zeros_like(img)
        img = cv2.resize(img, (int(len(img[0]) / zoomlevel), int(len(img) /
                                                                 zoomlevel)), interpolation = cv2.INTER_AREA)
        dst[0:len(img), 0:len(img[0])] = img
        return dst

    def zoomdataout(data, level):
        zoomdata = copy.deepcopy(data)
        zoomdata[0] = (zoomdata[0][0], zoomdata[0][1]/level,
                        zoomdata[0][2]/level, zoomdata[0][3]/level,
                       zoomdata[0][4]/level)
        return zoomdata
    apply_transform(zoomout, "zoomout", [1.4, 1.8, 2.2, 2.6, 3],
                    zoomdataout)

    # upsampling
    apply_transform(lambda image, samplelevel: cv2.resize(img, (int(len(img[0]) * samplelevel), int(len(img) * samplelevel)),
                                                    interpolation =
                                                    cv2.INTER_AREA),
                    "upsample", [2, 3, 4])

df = pd.DataFrame(names, columns=["colummn"])
df.to_csv("lists.txt", index=False, header=False)

master.to_csv("newlabels.csv", index=False)
