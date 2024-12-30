QT       += core gui widgets webenginewidgets webenginecore

TEMPLATE = app
TARGET = {name}
DESTDIR = dist

OBJECTS_DIR = build/obj
MOC_DIR = build/moc
RCC_DIR = build/rcc
UI_DIR = build/ui

macx {
    QMAKE_MACOSX_DEPLOYMENT_TARGET = 13.0
    QMAKE_MAC_SDK = macosx13.0
    QMAKE_INFO_PLIST = $${PWD}/Info.plist
    QMAKE_APPLE_DEVICE_ARCHS = x86_64 arm64
    ;ICON
}

win32{
    ;RC_ICONS
}

SOURCES += \
    browser.cpp \
    main.cpp \
    webinterceptor.cpp

HEADERS += \
    browser.h \
    webinterceptor.h \
    opt.h

