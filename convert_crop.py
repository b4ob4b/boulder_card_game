# convert and crop
# https://github.com/jarrekk/imgkit
import imgkit
import os
import sys

path = sys.argv[1]

options = {
    'format': 'png',
    'crop-h': '436',
    'crop-w': '300',
    'crop-x': '362',
    'crop-y': '10'
}

files_to_convert = []
for file in os.listdir(path + "/output_html/"):
    if file.endswith(".html"):
    	files_to_convert.append(os.path.join(path + "/output_html/", file))


for file in files_to_convert:
	imgkit.from_file(
		file,
		file.replace('html','png'), 
		options = options
	)

print("\n>> Cards were converted successfully!")





