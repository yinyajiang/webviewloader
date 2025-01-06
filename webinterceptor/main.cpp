#include <QApplication>
#include <QCommandLineParser>
#include <QFileInfo>
#include <QIcon>
#include "browser.h"
#include <QProcessEnvironment>
#ifdef _WIN32
    #include <io.h>
    #include <fcntl.h>
#endif
#include <QNetworkProxyFactory>

QString decodeMaybeBase64(QString base64);

int main(int argc, char *argv[]) {
    qunsetenv("LANG");
    #ifdef _WIN32
        _setmode(_fileno(stdout), _O_BINARY);
    #endif

    QCoreApplication::setAttribute(Qt::AA_EnableHighDpiScaling);

    QApplication app(argc, argv);

    QCommandLineParser parser;
    parser.setApplicationDescription("Playable url interceptor");
    parser.addHelpOption();
    parser.addPositionalArgument("url", "Target URL");

    //name, description, valueName, defaultValue
    parser.addOption(QCommandLineOption("title", "Window title", "title"));
    parser.addOption(QCommandLineOption("banner", "Banner text", "banner", "Please start playing the video first, and then it will start parsing."));
    parser.addOption(QCommandLineOption("banner-color", "Banner background color", "color", "#FF4E50"));
    parser.addOption(QCommandLineOption("ua", "User agent", "ua"));
    parser.addOption(QCommandLineOption("width", "Window width", "width", "1024"));
    parser.addOption(QCommandLineOption("height", "Window height", "height", "768"));
    parser.addOption(QCommandLineOption("icon", "Window icon", "icon", "icon.ico"));
    parser.addOption(QCommandLineOption("address", "Show address bar"));
    parser.addOption(QCommandLineOption("banner-font-color", "Banner font color", "color", "#F9F9F9"));
    parser.addOption(QCommandLineOption("win-color", "Windows color", "color", "#121212"));
    parser.addOption(QCommandLineOption("forever", "Run forever"));
    parser.addOption(QCommandLineOption("system-proxy", "Use system proxy"));
    parser.process(app);

    const QStringList args = parser.positionalArguments();
    if (args.isEmpty()) {
        parser.showHelp(1);
        return 1;
    }
    QString url = args.first();

    QString title = parser.value("title");
    if (title.isEmpty()) {
        QFileInfo fi(argv[0]);
        title = fi.baseName();
    }

    QString iconPath = parser.value("icon");
    if (QFile::exists(iconPath)) {
        app.setWindowIcon(QIcon(iconPath));
    }

    Options opt;
    opt.url = url;
    opt.title = decodeMaybeBase64(title);
    opt.ua = parser.value("ua");
    opt.width = parser.value("width").toInt();
    opt.height = parser.value("height").toInt();
    opt.banner = decodeMaybeBase64(parser.value("banner"));
    opt.bannerColor = parser.value("banner-color");
    opt.bannerFontColor = parser.value("banner-font-color");
    opt.showAddress = parser.isSet("address");
    opt.winColor = parser.value("win-color");
    opt.isforever = parser.isSet("forever");

    if(parser.isSet("system-proxy")) {
        QList<QNetworkProxy> proxyList = QNetworkProxyFactory::systemProxyForQuery();
        if(!proxyList.isEmpty()) {
            QNetworkProxy::setApplicationProxy(proxyList.first());
        }
    }

    Browser browser(opt);
    browser.show();
    return app.exec();
}


QString decodeMaybeBase64(QString base64) {
    if(base64.startsWith("base64://")) {
        base64 = base64.mid(QString("base64://").length());
    }else{
        return base64;
    }
    QByteArray decoded = QByteArray::fromBase64(base64.toUtf8());
    return QString::fromUtf8(decoded);
}
