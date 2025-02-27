#!/bin/bash

# 指定文件夹路径
folder_path="./"  # 喵！主人要替换成自己的文件夹路径喵！

# 遍历文件夹下的所有mp4文件
find "$folder_path" -name "*.mp4" -print0 | while IFS= read -r -d $'\0' file; do
  # 获取文件名（不包含路径）
  filename=$(basename "$file")

  # 输出正在处理的文件名
  echo "正在处理: $filename"

  # 使用ffmpeg添加视频流标签
  ffmpeg -i "$file" -c:v copy -c:a copy -tag:v hvc1 "${file%.mp4}_tagged.mp4"

  # 喵！主人可以根据需要选择是否覆盖原文件喵！
  # 如果主人想覆盖原文件，可以取消注释下面这行喵！
  # mv "${file%.mp4}_tagged.mp4" "$file"

  echo "处理完成: $filename"
done

echo "全部处理完成喵！"

