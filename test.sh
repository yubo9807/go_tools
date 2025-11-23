FILE="practical"

if [ -d "$FILE" ]; then
  cd $FILE
  git pull origin main
else
  git clone https://github.com/yubo9807/practical &&
  cd $FILE
fi

npm install &&
npm run build &&
rm -rf node_modules
