# UTTC_hackathon_back
ハッカソン用のバックエンドリポジトリ

tmp:ローカルデータベースを立ち上げるためのdockerディレクトリ

controller, dao, usecaseはエンドポイントごとにコードを分割。
daoでは加えてdao/init_dao.goでdaoアクセスを管理