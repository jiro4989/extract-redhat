# extract-redhat
RedHatのセキュリティアップデート情報のページのHTMLから覚寿情報を取得してタブ区切
りのCSVファイルとして出力します。

## 使い方
RedHatEnterpriseEditionのセキュリティアップデートのページを開き、ブラウザのF12か
ら開発者画面を起動する。

HTMLを表示し、ルート要素をinlin copyする。

data/target.htmlの中身に貼り付ける。

端末から実行ファイルを動かす。コマンドラインオプションは指定せずとも動くはず。

コマンドラインオプションからHTMLファイルを指定することも可能。  
詳細は`extract-redhat.exe -h`を参照。
