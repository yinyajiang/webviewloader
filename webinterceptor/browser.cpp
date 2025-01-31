#include <iostream>
#include "browser.h"
#include <QVBoxLayout>
#include <QHBoxLayout>
#include <QPushButton>
#include <QJsonObject>
#include <QJsonDocument>
#include <QCloseEvent>
#include "webinterceptor.h"


Browser::Browser(Options opt)
    : QMainWindow() {
    setWindowTitle(opt.title);

    QWidget* centralWidget = new QWidget(this);
    if(!opt.winColor.isEmpty() && opt.winColor != "none" && opt.winColor != "null"){
        centralWidget->setStyleSheet(QString("QWidget { background-color: %1;}").arg(opt.winColor));
    }

    setCentralWidget(centralWidget);

    QVBoxLayout* layout = new QVBoxLayout(centralWidget);

    if (opt.showAddress) {
        QHBoxLayout* addressLayout = new QHBoxLayout();
        m_urlEdit = new QLineEdit(this);
        m_urlEdit->setText(opt.url);
        m_urlEdit->setStyleSheet(QString(
                                   "QLineEdit {"
                                   "    padding: 6px 20px;"
                                   "    font-size: 13px;"
                                   "    border-radius: 3px;"
                                   "    border: 1px solid %1;"
                                   "    margin-right: 5px;"
                                   "}").arg(opt.bannerColor));

        QPushButton* loadButton = new QPushButton("Go", this);
        loadButton->setStyleSheet(QString(
                                      "QPushButton {"
                                      "    padding: 6px 20px;"
                                      "    font-size: 13px;"
                                      "    border-radius: 3px;"
                                      "    background-color: %1;"
                                      "    color: white;"
                                      "    border: none;"
                                      "}"
                                      "QPushButton:hover {"
                                      "    background-color: #ff6668;"
                                      "}").arg(opt.bannerColor));

        connect(loadButton, &QPushButton::clicked, this, &Browser::loadUrl);

        addressLayout->addWidget(m_urlEdit);
        addressLayout->addWidget(loadButton);
        addressLayout->setContentsMargins(15, 8, 15, 0);
        layout->addLayout(addressLayout);
    }

    QString bannerText = QString("%1").arg(opt.banner);
    QLabel* bannerLabel = new QLabel(bannerText, this);
    bannerLabel->setAlignment(Qt::AlignCenter);
    bannerLabel->setStyleSheet(QString(
                                   "QLabel {"
                                   "    background-color: %1;"
                                   "    color: %2;"
                                   "    padding: 6px 20px;"
                                   "    font-size: 14px;"
                                   "    font-weight: bold;"
                                   "    border-radius: 3px;"
                                   "    margin: 8px 15px;"
                                   "    max-height: 22px;"
                                   "    letter-spacing: 1px;"
                                   "}"
                                   "QLabel::first-letter {"
                                   "    font-size: 14px;"
                                   "    margin-right: 8px;"
                                     "}").arg(opt.bannerColor).arg(opt.bannerFontColor));
    bannerLabel->setMaximumHeight(40);

    layout->addWidget(bannerLabel);

    m_webView = new QWebEngineView(this);
    m_profile = new WebInterceptor(m_webView, opt);
    m_page = new QWebEnginePage(m_profile, this);
    connect(m_page, &QWebEnginePage::urlChanged, this, &Browser::urlChanged);
    m_webView->setPage(m_page);

    layout->addWidget(m_webView);
    layout->setContentsMargins(0, 0, 0, 0);
    layout->setSpacing(0);

    resize(opt.width, opt.height);
    m_webView->load(QUrl(opt.url));
}

Browser::~Browser() {
    m_webView->setPage(nullptr);
    delete m_page;
}

void Browser::closeEvent(QCloseEvent* event) {
    QJsonObject obj;
    obj["error"] = "Closed by user.";
    QJsonDocument doc(obj);
    std::cout << doc.toJson(QJsonDocument::Compact).toStdString() << std::endl;
    std::cout.flush();
    QMainWindow::closeEvent(event);
}

void Browser::loadUrl() {
    if (m_urlEdit) {
        m_webView->load(QUrl(m_urlEdit->text()));
    }
}

void Browser::urlChanged(const QUrl& url) {
    if (m_urlEdit) {
        m_urlEdit->setText(url.toString());
    }
}
