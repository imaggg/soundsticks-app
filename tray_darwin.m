#import <Cocoa/Cocoa.h>

// Borderless NSWindow returns NO for canBecomeKeyWindow by default,
// which prevents keyboard input in WKWebView. Override for all windows.
@interface NSWindow (SSKeyable)
@end
@implementation NSWindow (SSKeyable)
- (BOOL)canBecomeKeyWindow { return YES; }
@end

extern void trayQuit(void);
extern void trayReconnect(void);
extern void trayConnectIP(const char *ip);
void showConnectIPPrompt(void);
extern void trayStatus(void);
extern void trayInfo(void);
extern void trayForget(void);
extern void trayForgetConfirmed(void);
extern void traySetLanguage(const char *lang);
extern void traySetTheme(const char *theme);
extern void trayToggleAutostart(void);
extern void traySetLedKeepAlive(int minutes);
extern void trayPopoverWillShow(void);
extern void trayPopoverDidHide(void);

@class SSTrayDelegate;
static SSTrayDelegate *_delegate;
static NSStatusItem   *_item;
static NSWindow       *_window;
static NSString *_currentLang  = @"en";
static NSString *_currentTheme = @"auto";
static BOOL      _autostartEnabled = NO;
static int       _ledKeepAlive = 0;
static NSString *_deviceIP = @"";

static NSString *ms(NSString *key) {
    static NSDictionary *table = nil;
    if (!table) table = @{
        @"status":         @{@"en":@"Status",              @"es":@"Estado",             @"fr":@"Statut",              @"uk":@"Статус"},
        @"info":           @{@"en":@"Info",                @"es":@"Info",               @"fr":@"Info",                @"uk":@"Інфо"},
        @"language":       @{@"en":@"Language",            @"es":@"Idioma",             @"fr":@"Langue",              @"uk":@"Мова"},
        @"reconnect":      @{@"en":@"Reconnect",           @"es":@"Reconectar",         @"fr":@"Reconnecter",         @"uk":@"Перепідключитись"},
        @"connect_ip":     @{@"en":@"Connect to IP…",      @"es":@"Conectar a IP…",     @"fr":@"Connexion par IP…",   @"uk":@"Підключитись за IP…"},
        @"connect_ip.title":@{@"en":@"Connect to IP",      @"es":@"Conectar a IP",      @"fr":@"Connexion par IP",    @"uk":@"Підключитись за IP"},
        @"connect_ip.msg": @{@"en":@"Enter the speaker's IP address on your network.",
                             @"es":@"Introduce la dirección IP del altavoz en tu red.",
                             @"fr":@"Saisissez l'adresse IP de l'enceinte sur votre réseau.",
                             @"uk":@"Введіть IP-адресу колонки у вашій мережі."},
        @"connect_ip.confirm":@{@"en":@"Connect",          @"es":@"Conectar",           @"fr":@"Connexion",           @"uk":@"Підключитись"},
        @"forget":         @{@"en":@"Forget Speaker", @"es":@"Olvidar Certificado",@"fr":@"Oublier le certificat",@"uk":@"Забути колонку"},
        @"quit":           @{@"en":@"Quit",    @"es":@"Salir",@"fr":@"Quitter ",@"uk":@"Вийти"},
        @"forget.title":   @{@"en":@"Forget Speaker?", @"es":@"¿Olvidar Certificado?",@"fr":@"Oublier le certificat ?",@"uk":@"Забути колонку?"},
        @"forget.msg":     @{@"en":@"App will remove the stored certificate and return to the setup screen. You can re-import it at any time.",
                             @"es":@"App eliminará el certificado guardado y volverá a la pantalla de configuración. Puedes volver a importarlo en cualquier momento.",
                             @"fr":@"App supprimera le certificat stocké et reviendra à l'écran de configuration. Vous pouvez le réimporter à tout moment.",
                             @"uk":@"Додаток видалить збережений сертифікат і повернеться до екрану налаштування."},
        @"forget.confirm": @{@"en":@"Forget",   @"es":@"Olvidar",  @"fr":@"Oublier",  @"uk":@"Забути"},
        @"cancel":         @{@"en":@"Cancel",          @"es":@"Cancelar",         @"fr":@"Annuler",             @"uk":@"Скасувати"},
        @"autostart":      @{@"en":@"Launch at Login",    @"es":@"Iniciar al arrancar",   @"fr":@"Lancer au démarrage",  @"uk":@"Запускати при вході"},
        @"led_sleep":      @{@"en":@"Led Sleep Override", @"es":@"Mantener LED activo",   @"fr":@"Maintenir LED actif",  @"uk":@"Не вимикати підсвітку"},
        @"led_off":        @{@"en":@"Off (Default)",      @"es":@"Apagado (Predeterminado)",@"fr":@"Désactivé (par défaut)",@"uk":@"Вимкнено (за замовчуванням)"},
        @"theme":          @{@"en":@"Theme",              @"es":@"Tema",               @"fr":@"Thème",               @"uk":@"Тема"},
        @"theme.light":    @{@"en":@"Light",              @"es":@"Claro",              @"fr":@"Clair",               @"uk":@"Світла"},
        @"theme.dark":     @{@"en":@"Dark",               @"es":@"Oscuro",             @"fr":@"Sombre",              @"uk":@"Темна"},
        @"theme.auto":     @{@"en":@"Auto (System)",      @"es":@"Auto (Sistema)",     @"fr":@"Auto (Système)",      @"uk":@"Авто (Система)"},
    };
    NSDictionary *row = table[key];
    return row[_currentLang] ?: row[@"en"] ?: key;
}

@interface SSTrayDelegate : NSObject
@end

@implementation SSTrayDelegate
- (void)doQuit:(id)sender      { trayQuit(); }
- (void)doReconnect:(id)sender { trayReconnect(); }
- (void)doConnectIP:(id)sender { showConnectIPPrompt(); }
- (void)doStatus:(id)sender    { trayStatus(); }
- (void)doInfo:(id)sender      { trayInfo(); }
- (void)doForget:(id)sender    { trayForget(); }
- (void)setLangEn:(id)sender      { traySetLanguage("en"); }
- (void)setLangEs:(id)sender      { traySetLanguage("es"); }
- (void)setLangFr:(id)sender      { traySetLanguage("fr"); }
- (void)setLangUk:(id)sender      { traySetLanguage("uk"); }
- (void)doAutostart:(id)sender    { trayToggleAutostart(); }
- (void)setThemeLight:(id)sender  { traySetTheme("light"); }
- (void)setThemeDark:(id)sender   { traySetTheme("dark"); }
- (void)setThemeAuto:(id)sender   { traySetTheme("auto"); }
- (void)setLedOff:(id)sender      { traySetLedKeepAlive(0); }
- (void)setLed60:(id)sender       { traySetLedKeepAlive(60); }
- (void)setLed120:(id)sender      { traySetLedKeepAlive(120); }

- (void)showPopover {
    if (!_window) return;
    NSRect btn = [_item.button.window convertRectToScreen:_item.button.frame];
    NSRect win = _window.frame;
    CGFloat x = NSMidX(btn) - win.size.width / 2;
    CGFloat y = NSMinY(btn) - win.size.height - 6;
    NSScreen *screen = _item.button.window.screen ?: [NSScreen mainScreen];
    NSRect    vis    = screen.visibleFrame;
    if (x < NSMinX(vis) + 4)                  x = NSMinX(vis) + 4;
    if (x + win.size.width > NSMaxX(vis) - 4) x = NSMaxX(vis) - 4 - win.size.width;
    [_window setFrameOrigin:NSMakePoint(x, y)];
    _window.alphaValue = 0;
    trayPopoverWillShow();
    [NSApp activateIgnoringOtherApps:YES];
    [_window makeKeyAndOrderFront:nil];
    [NSAnimationContext runAnimationGroup:^(NSAnimationContext *ctx) {
        ctx.duration = 0.15;
        [_window.animator setAlphaValue:1.0];
    }];
}

- (void)hidePopover {
    if (!_window.isVisible) return;
    trayPopoverDidHide();
    [NSAnimationContext runAnimationGroup:^(NSAnimationContext *ctx) {
        ctx.duration = 0.10;
        [_window.animator setAlphaValue:0.0];
    } completionHandler:^{
        [_window orderOut:nil];
        _window.alphaValue = 1.0;
    }];
}

- (void)handleClick:(id)sender {
    NSEvent *e = [NSApp currentEvent];
    if (e.type == NSEventTypeRightMouseUp ||
        (e.modifierFlags & NSEventModifierFlagControl)) {
        NSMenu *menu = [[NSMenu alloc] init];

        NSMenuItem *st = [[NSMenuItem alloc] initWithTitle:ms(@"status")
            action:@selector(doStatus:) keyEquivalent:@""];
        st.target = self; [menu addItem:st];

        NSMenuItem *inf = [[NSMenuItem alloc] initWithTitle:ms(@"info")
            action:@selector(doInfo:) keyEquivalent:@""];
        inf.target = self; [menu addItem:inf];

        [menu addItem:[NSMenuItem separatorItem]];

        // Language submenu
        NSMenuItem *langItem = [[NSMenuItem alloc] initWithTitle:ms(@"language")
            action:nil keyEquivalent:@""];
        NSMenu *langMenu = [[NSMenu alloc] init];
        struct { NSString *title; NSString *code; SEL sel; } langs[] = {
            {@"English",    @"en", @selector(setLangEn:)},
            {@"Español",    @"es", @selector(setLangEs:)},
            {@"Français",   @"fr", @selector(setLangFr:)},
            {@"Українська", @"uk", @selector(setLangUk:)},
        };
        for (int i = 0; i < 4; i++) {
            NSMenuItem *li = [[NSMenuItem alloc] initWithTitle:langs[i].title
                action:langs[i].sel keyEquivalent:@""];
            li.target = self;
            li.state = [langs[i].code isEqualToString:_currentLang]
                ? NSControlStateValueOn : NSControlStateValueOff;
            [langMenu addItem:li];
        }
        langItem.submenu = langMenu;
        [menu addItem:langItem];

        // Theme submenu
        NSMenuItem *themeItem = [[NSMenuItem alloc] initWithTitle:ms(@"theme")
            action:nil keyEquivalent:@""];
        NSMenu *themeMenu = [[NSMenu alloc] init];
        struct { NSString *title; NSString *code; SEL sel; } themes[] = {
            {ms(@"theme.light"), @"light", @selector(setThemeLight:)},
            {ms(@"theme.dark"),  @"dark",  @selector(setThemeDark:)},
            {ms(@"theme.auto"),  @"auto",  @selector(setThemeAuto:)},
        };
        for (int i = 0; i < 3; i++) {
            NSMenuItem *ti = [[NSMenuItem alloc] initWithTitle:themes[i].title
                action:themes[i].sel keyEquivalent:@""];
            ti.target = self;
            ti.state = [themes[i].code isEqualToString:_currentTheme]
                ? NSControlStateValueOn : NSControlStateValueOff;
            [themeMenu addItem:ti];
        }
        themeItem.submenu = themeMenu;
        [menu addItem:themeItem];

        NSMenuItem *ast = [[NSMenuItem alloc] initWithTitle:ms(@"autostart")
            action:@selector(doAutostart:) keyEquivalent:@""];
        ast.target = self;
        ast.state = _autostartEnabled ? NSControlStateValueOn : NSControlStateValueOff;
        [menu addItem:ast];

        NSMenuItem *ledItem = [[NSMenuItem alloc] initWithTitle:ms(@"led_sleep")
            action:nil keyEquivalent:@""];
        NSMenu *ledMenu = [[NSMenu alloc] init];

        NSMenuItem *led0 = [[NSMenuItem alloc] initWithTitle:ms(@"led_off")
            action:@selector(setLedOff:) keyEquivalent:@""];
        led0.target = self;
        led0.state = (_ledKeepAlive == 0) ? NSControlStateValueOn : NSControlStateValueOff;
        [ledMenu addItem:led0];

        NSMenuItem *led60 = [[NSMenuItem alloc] initWithTitle:@"60 min"
            action:@selector(setLed60:) keyEquivalent:@""];
        led60.target = self;
        led60.state = (_ledKeepAlive == 60) ? NSControlStateValueOn : NSControlStateValueOff;
        [ledMenu addItem:led60];

        NSMenuItem *led120 = [[NSMenuItem alloc] initWithTitle:@"120 min"
            action:@selector(setLed120:) keyEquivalent:@""];
        led120.target = self;
        led120.state = (_ledKeepAlive == 120) ? NSControlStateValueOn : NSControlStateValueOff;
        [ledMenu addItem:led120];

        ledItem.submenu = ledMenu;
        [menu addItem:ledItem];

        [menu addItem:[NSMenuItem separatorItem]];

        NSMenuItem *r = [[NSMenuItem alloc] initWithTitle:ms(@"reconnect")
            action:@selector(doReconnect:) keyEquivalent:@""];
        r.target = self; [menu addItem:r];

        NSMenuItem *cip = [[NSMenuItem alloc] initWithTitle:ms(@"connect_ip")
            action:@selector(doConnectIP:) keyEquivalent:@""];
        cip.target = self; [menu addItem:cip];

        [menu addItem:[NSMenuItem separatorItem]];

        NSMenuItem *fg = [[NSMenuItem alloc] initWithTitle:ms(@"forget")
            action:@selector(doForget:) keyEquivalent:@""];
        fg.target = self; [menu addItem:fg];

        [menu addItem:[NSMenuItem separatorItem]];

        NSMenuItem *q = [[NSMenuItem alloc] initWithTitle:ms(@"quit")
            action:@selector(doQuit:) keyEquivalent:@"q"];
        q.target = self; [menu addItem:q];
        NSPoint loc = NSMakePoint(0, _item.button.bounds.size.height);
        [menu popUpMenuPositioningItem:nil atLocation:loc inView:_item.button];
        return;
    }
    if (_window.isVisible) {
        [self hidePopover];
    } else {
        [self showPopover];
    }
}
@end

void setupTray(void) {
    [NSApplication sharedApplication];
    [NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];

    _delegate = [[SSTrayDelegate alloc] init];

    NSStatusBar *bar = [NSStatusBar systemStatusBar];
    _item = [bar statusItemWithLength:NSVariableStatusItemLength];

    NSImage *img = nil;
    if (@available(macOS 11.0, *)) {
        img = [NSImage imageWithSystemSymbolName:@"speaker.wave.2.fill"
                          accessibilityDescription:nil];
        [img setTemplate:YES];
    }
    if (img) { _item.button.image = img; } else { _item.button.title = @"SS"; }

    _item.button.target = _delegate;
    _item.button.action = @selector(handleClick:);
    [_item.button sendActionOn:NSEventMaskLeftMouseUp | NSEventMaskRightMouseUp];
}

void attachPopoverWindow(void *nsWindowPtr) {
    if (!nsWindowPtr) return;
    _window = (__bridge NSWindow *)nsWindowPtr;

    // Re-apply after webview.New() which may have temporarily activated the app.
    [NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];

    _window.styleMask          = NSWindowStyleMaskBorderless;
    _window.hasShadow          = YES;
    _window.level              = NSStatusWindowLevel;
    _window.collectionBehavior = NSWindowCollectionBehaviorCanJoinAllSpaces |
                                 NSWindowCollectionBehaviorTransient        |
                                 NSWindowCollectionBehaviorIgnoresCycle;
    _window.movableByWindowBackground = NO;
    _window.hidesOnDeactivate         = NO;

    NSView *content = _window.contentView;
    content.wantsLayer          = YES;
    content.layer.cornerRadius  = 10;
    content.layer.masksToBounds = YES;

    [NSEvent addGlobalMonitorForEventsMatchingMask:NSEventMaskLeftMouseDown | NSEventMaskRightMouseDown
        handler:^(NSEvent *e) {
            if (_window.isVisible) [_delegate hidePopover];
        }];

    [_window orderOut:nil];
}

// ── Status alert ──────────────────────────────────────────────────────────────

void showStatusAlert(const char *ctext) {
    NSString *text = [NSString stringWithUTF8String:ctext];
    dispatch_async(dispatch_get_main_queue(), ^{
        NSAlert *a = [[NSAlert alloc] init];
        a.messageText     = @"SoundSticks Status";
        a.informativeText = text;
        [a addButtonWithTitle:@"OK"];
        [a runModal];
    });
}

// ── Info panel ("About" style) ────────────────────────────────────────────────

void showInfoPanel(const char *cname, const char *cfirmware,
                   const char *cserial, const char *cmac, const char *cuuid) {
    NSString *name     = [NSString stringWithUTF8String:cname];
    NSString *firmware = [NSString stringWithUTF8String:cfirmware];
    NSString *serial   = [NSString stringWithUTF8String:cserial];
    NSString *mac      = [NSString stringWithUTF8String:cmac];
    NSString *uuid     = [NSString stringWithUTF8String:cuuid];

    dispatch_async(dispatch_get_main_queue(), ^{
        NSPanel *panel = [[NSPanel alloc]
            initWithContentRect:NSMakeRect(0, 0, 340, 268)
            styleMask:NSWindowStyleMaskTitled | NSWindowStyleMaskClosable
            backing:NSBackingStoreBuffered
            defer:NO];
        panel.title = @"About SoundSticks";
        panel.level = NSFloatingWindowLevel;
        [panel center];

        NSView *cv = panel.contentView;

        // Speaker icon
        if (@available(macOS 11.0, *)) {
            NSImage *img = [NSImage imageWithSystemSymbolName:@"hifispeaker.fill"
                                      accessibilityDescription:nil];
            NSImageSymbolConfiguration *cfg = [NSImageSymbolConfiguration
                configurationWithPointSize:52 weight:NSFontWeightLight];
            NSImageView *iv = [NSImageView imageViewWithImage:
                [img imageWithSymbolConfiguration:cfg]];
            iv.frame = NSMakeRect(130, 196, 80, 56);
            iv.imageScaling = NSImageScaleProportionallyDown;
            [cv addSubview:iv];
        }

        // Device name
        NSTextField *nameLabel = [NSTextField labelWithString:name];
        nameLabel.font      = [NSFont boldSystemFontOfSize:16];
        nameLabel.alignment = NSTextAlignmentCenter;
        nameLabel.frame     = NSMakeRect(0, 164, 340, 24);
        [cv addSubview:nameLabel];

        // Separator
        NSBox *sep = [[NSBox alloc] initWithFrame:NSMakeRect(30, 150, 280, 1)];
        sep.boxType = NSBoxSeparator;
        [cv addSubview:sep];

        // Key-value rows
        NSArray *rows = @[
            @[@"Firmware", firmware],
            @[@"Serial",   serial  ],
            @[@"MAC",      mac     ],
            @[@"UUID",     uuid    ],
        ];
        CGFloat y = 124;
        for (NSArray *row in rows) {
            NSTextField *kl = [NSTextField labelWithString:row[0]];
            kl.font       = [NSFont systemFontOfSize:11];
            kl.textColor  = [NSColor secondaryLabelColor];
            kl.alignment  = NSTextAlignmentRight;
            kl.frame      = NSMakeRect(16, y, 80, 16);
            [cv addSubview:kl];

            NSTextField *vl = [NSTextField labelWithString:row[1]];
            vl.font  = [NSFont monospacedSystemFontOfSize:11 weight:NSFontWeightRegular];
            vl.frame = NSMakeRect(106, y, 218, 16);
            [cv addSubview:vl];
            y -= 26;
        }

        // Separator above credits
        NSBox *sep2 = [[NSBox alloc] initWithFrame:NSMakeRect(30, 28, 280, 1)];
        sep2.boxType = NSBoxSeparator;
        [cv addSubview:sep2];

        // "Made by imaggg with Claude" with clickable links
        NSMutableParagraphStyle *ps = [[NSMutableParagraphStyle alloc] init];
        ps.alignment = NSTextAlignmentCenter;

        NSMutableAttributedString *credit = [[NSMutableAttributedString alloc] init];
        NSDictionary *plain = @{NSFontAttributeName: [NSFont systemFontOfSize:10],
                                NSForegroundColorAttributeName: [NSColor secondaryLabelColor],
                                NSParagraphStyleAttributeName: ps};
        NSDictionary *link1 = @{NSFontAttributeName: [NSFont systemFontOfSize:10],
                                NSLinkAttributeName: [NSURL URLWithString:@"https://imaggg.com"],
                                NSForegroundColorAttributeName: [NSColor linkColor],
                                NSParagraphStyleAttributeName: ps};
        NSDictionary *link2 = @{NSFontAttributeName: [NSFont systemFontOfSize:10],
                                NSLinkAttributeName: [NSURL URLWithString:@"https://claude.ai"],
                                NSForegroundColorAttributeName: [NSColor linkColor],
                                NSParagraphStyleAttributeName: ps};
        [credit appendAttributedString:[[NSAttributedString alloc] initWithString:@"Made by " attributes:plain]];
        [credit appendAttributedString:[[NSAttributedString alloc] initWithString:@"imaggg" attributes:link1]];
        [credit appendAttributedString:[[NSAttributedString alloc] initWithString:@" with " attributes:plain]];
        [credit appendAttributedString:[[NSAttributedString alloc] initWithString:@"Claude" attributes:link2]];

        NSTextField *creditField = [[NSTextField alloc] initWithFrame:NSMakeRect(0, 6, 340, 18)];
        creditField.attributedStringValue = credit;
        creditField.bordered = NO;
        creditField.drawsBackground = NO;
        creditField.editable = NO;
        creditField.selectable = YES;
        creditField.allowsEditingTextAttributes = YES;
        [cv addSubview:creditField];

        [panel makeKeyAndOrderFront:nil];
        [NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];
        [NSApp activateIgnoringOtherApps:YES];
    });
}

// ── Forget confirmation ───────────────────────────────────────────────────────

void showForgetConfirm(void) {
    NSString *title   = ms(@"forget.title");
    NSString *msg     = ms(@"forget.msg");
    NSString *confirm = ms(@"forget.confirm");
    NSString *cancel  = ms(@"cancel");
    dispatch_async(dispatch_get_main_queue(), ^{
        NSAlert *a = [[NSAlert alloc] init];
        a.messageText     = title;
        a.informativeText = msg;
        a.alertStyle = NSAlertStyleWarning;
        [a addButtonWithTitle:confirm];
        [a addButtonWithTitle:cancel];
        if ([a runModal] == NSAlertFirstButtonReturn) {
            trayForgetConfirmed();
        }
    });
}

// ── Manual "Connect to IP" prompt ─────────────────────────────────────────────

void setDeviceIPField(const char *ip) {
    _deviceIP = ip ? [NSString stringWithUTF8String:ip] : @"";
}

void showConnectIPPrompt(void) {
    NSString *title   = ms(@"connect_ip.title");
    NSString *msg     = ms(@"connect_ip.msg");
    NSString *confirm = ms(@"connect_ip.confirm");
    NSString *cancel  = ms(@"cancel");
    NSString *prefill = _deviceIP ?: @"";
    dispatch_async(dispatch_get_main_queue(), ^{
        NSAlert *a = [[NSAlert alloc] init];
        a.messageText     = title;
        a.informativeText = msg;
        [a addButtonWithTitle:confirm];
        [a addButtonWithTitle:cancel];

        NSTextField *input = [[NSTextField alloc] initWithFrame:NSMakeRect(0, 0, 220, 24)];
        input.stringValue        = prefill;
        input.placeholderString  = @"192.168.1.95";
        a.accessoryView = input;
        [a.window setInitialFirstResponder:input];

        [NSApp activateIgnoringOtherApps:YES];
        if ([a runModal] == NSAlertFirstButtonReturn) {
            NSString *ip = [input.stringValue
                stringByTrimmingCharactersInSet:[NSCharacterSet whitespaceAndNewlineCharacterSet]];
            if (ip.length > 0) {
                trayConnectIP(ip.UTF8String);
            }
        }
    });
}

// ── Language checkmark ────────────────────────────────────────────────────────

void setLanguageMenuCheck(const char *lang) {
    _currentLang = [NSString stringWithUTF8String:lang];
}

// ── Autostart checkmark ───────────────────────────────────────────────────────

void setAutostartCheck(int enabled) {
    _autostartEnabled = enabled != 0;
}

void setLedKeepAliveCheck(int minutes) {
    _ledKeepAlive = minutes;
}

void setThemeMenuCheck(const char *theme) {
    _currentTheme = [NSString stringWithUTF8String:theme];
}

// ── Window height ─────────────────────────────────────────────────────────────

void setWindowHeight(int h) {
    dispatch_async(dispatch_get_main_queue(), ^{
        if (!_window) return;
        NSScreen *screen = _window.screen ?: [NSScreen mainScreen];
        CGFloat maxH = screen.visibleFrame.size.height * 0.9;
        CGFloat newH = MIN((CGFloat)h, maxH);
        if (newH < 50) return;
        NSRect f = _window.frame;
        f.origin.y += f.size.height - newH;
        f.size.height = newH;
        [_window setFrame:f display:NO];
    });
}
