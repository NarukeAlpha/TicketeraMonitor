package Monitor

import (
	"log"
	"math/rand"
	"time"

	"github.com/playwright-community/playwright-go"
)

func TaskInit(pL []ProxyStruct, wbKey string) {

	for {
		for _, proxy := range pL {
			Task(proxy, wbKey)
			if err := recover(); err != nil {
				log.Printf("Recovered from panic: %v", err)
			}
			time.Sleep(3 * time.Minute)

		}
		time.Sleep(10 * time.Minute)
	}

}

func PlaywrightInit(proxy ProxyStruct) playwright.BrowserContext {
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}
	device := pw.Devices[IphoneUserAgentList[rand.Intn(len(IphoneUserAgentList)-1)]]
	pwProxyStrct := playwright.Proxy{
		Server:   proxy.ip,
		Username: &proxy.usr,
		Password: &proxy.pw,
	}

	browser, err := pw.WebKit.LaunchPersistentContext("", playwright.BrowserTypeLaunchPersistentContextOptions{
		Viewport:          device.Viewport,
		UserAgent:         playwright.String(device.UserAgent),
		DeviceScaleFactor: playwright.Float(device.DeviceScaleFactor),
		IsMobile:          playwright.Bool(device.IsMobile),
		HasTouch:          playwright.Bool(device.HasTouch),
		Headless:          playwright.Bool(false),
		ColorScheme:       playwright.ColorSchemeDark,
		Proxy:             &pwProxyStrct,
		IgnoreDefaultArgs: []string{
			"--enable-automation",
		},
	})
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	return browser
}

func Task(proxy ProxyStruct, wbKey string) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Recovered from panic: %v assuming bad proxy", err)
		}
	}()

	browser := PlaywrightInit(proxy)

	script := playwright.Script{
		Content: playwright.String(`
    const defaultGetter = Object.getOwnPropertyDescriptor(
      Navigator.prototype,
      "webdriver"
    ).get;
    defaultGetter.apply(navigator);
    defaultGetter.toString();
    Object.defineProperty(Navigator.prototype, "webdriver", {
      set: undefined,
      enumerable: true,
      configurable: true,
      get: new Proxy(defaultGetter, {
        apply: (target, thisArg, args) => {
          Reflect.apply(target, thisArg, args);
          return false;
        },
      }),
    });
    const patchedGetter = Object.getOwnPropertyDescriptor(
      Navigator.prototype,
      "webdriver"
    ).get;
    patchedGetter.apply(navigator);
    patchedGetter.toString();
  `),
	}
	err := browser.AddInitScript(script)
	if err != nil {
		log.Fatalf("could not add initialization script: %v", err)
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)

	}
	if _, err := page.Goto("https://www.whatsmyip.org/"); err != nil {
		log.Panicf("could not goto: %v", err)

	}

	page2, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)

	}
	if _, err = page2.Goto("https://bot.sannysoft.com/"); err != nil {
		log.Panicf("could not goto: %v", err)
	}
	AssertErrorToNil("Failed to close Page", page.Close())
	AssertErrorToNil("Failed to close Page", page2.Close())

	page3, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)

	}
	if _, err = page3.Goto("https://cpr.ticketera.com/tickets/series/badbunnymostwantedtour"); err != nil {
		log.Panicf("could not goto: %v", err)
	}
	//err = page3.Locator("#challenge-stage").WaitFor(playwright.LocatorWaitForOptions{
	//	State:   playwright.WaitForSelectorStateVisible,
	//	Timeout: playwright.Float(10000),
	//})
	//if err != nil {
	//	log.Panicf("error waiting for selector: %v", err)
	//}

	time.Sleep(30 * time.Second)

	if page3.URL() == "https://cpr.ticketera.com/tickets/series/badbunnymostwantedtour" {
		locator := page3.Locator("#log-in-here")
		exists, err := locator.Count()
		if err != nil {
			log.Panicf("could not count element: at log in here %v", err)
		}
		if exists > 0 {
			log.Printf("Event is open")
			//send webhook
			AssertErrorToNil("Failed to close Page", page3.Close())
			AssertErrorToNil("Failed to close Browser", browser.Close())

		}
		locator = page3.Locator("#challenge-stage")
		exists, err = locator.Count()
		if err != nil {
			log.Panicf("could not count element: at log in here %v", err)
		}
		if exists > 0 {
			log.Printf("Captcha is cycling, trying with new proxy")
			AssertErrorToNil("Failed to close Page", page3.Close())
			AssertErrorToNil("Failed to close Browser", browser.Close())
			return

		}

	} else if page3.URL() == "https://tixtrack.queue-it.net/afterevent.aspx?c=tixtrack&e=badbunnymostwanted&t=https%3A%2F%2Fcpr.ticketera.com%2Ftickets%2Fseries%2Fbadbunnymostwantedtour&cid=en-US" {
		log.Printf("Event is still closed")
	} else {
		log.Printf("somethign went wrong")
	}
	//send webhook

	time.Sleep(3 * time.Minute)
	log.Printf("finished task for proxy :%v", proxy.ip)
	AssertErrorToNil("Failed to close Page", page3.Close())
	AssertErrorToNil("Failed to close Browser", browser.Close())
}
