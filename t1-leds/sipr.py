#!/usr/bin/env python3
import sounddevice
import numpy
import os
import sys

leds = {
    "red": 1,
    "green": 2
}

def print_sound(indata, outdata, frames, time, status):
    volume_norm = numpy.linalg.norm(indata)*10
    if volume_norm > 14:
        print(volume_norm)
        if volume_norm < 25:
            print(f"amp: {int(volume_norm)}, LED: green ")
            os.system(f"echo {leds['green']} | tee {devName}")
            print()
        elif volume_norm > 30:
            print(f"amp: {int(volume_norm)}, LED: red ")
            os.system(f"echo -n {leds['red']} | tee {devName}")
            print()


if __name__ == "__main__":
    rtd = sys.argv[1]
    devName = sys.argv[2]

    with sounddevice.Stream(callback=print_sound):
        sounddevice.sleep(int(rtd)*1000)
