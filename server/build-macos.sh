#!/bin/sh

APP="初始化工具.app"
mkdir -p $APP/Contents/{MacOS,Resources}
go build -o $APP/Contents/MacOS/gin-gf-vue-admin
cat > $APP/Contents/Info.plist << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleExecutable</key>
	<string>gin-gf-vue-admin</string>
	<key>CFBundleIconFile</key>
	<string>icon.icns</string>
	<key>CFBundleIdentifier</key>
	<string>com.gin.vue.admin</string>
</dict>
</plist>
EOF
cp icons/icon.icns $APP/Contents/Resources/icon.icns
find $APP