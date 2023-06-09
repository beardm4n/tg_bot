from pydub import AudioSegment
from pathlib import Path
import sys

# Получаем аргументы, переданные в командной строке 
audioName = sys.argv[1:][0]
inputFilePath = sys.argv[1:][1]
outputFilePath = sys.argv[1:][2]

input_file = Path(inputFilePath + '/' + audioName + '.mp3')
output_file = Path(outputFilePath + '/' + audioName + '.oga')

def convert_to_oga(input_file, output_file):
    audio = AudioSegment.from_mp3(input_file)
    audio.export(output_file, format='ogg')

    return ''

convert_to_oga(input_file, output_file)