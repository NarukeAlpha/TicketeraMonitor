package Monitor

import (
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/playwright-community/playwright-go"
)

func TaskInit(pL []ProxyStruct, wbKey string) {
	errch := make(chan error)
	var err error
	for {
		for _, proxy := range pL {
			go Task(proxy, wbKey, errch)
			err = <-errch
			if err != nil {
				continue
			} else {
				time.Sleep(2 * time.Minute)
			}
		}
	}

}

func PlaywrightInit(proxy ProxyStruct, pw *playwright.Playwright) (playwright.BrowserContext, error) {

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
		//	Headless:          playwright.Bool(false),
		ColorScheme: playwright.ColorSchemeDark,
		Proxy:       &pwProxyStrct,
		IgnoreDefaultArgs: []string{
			"--enable-automation",
		},
	})
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}

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
	err = browser.AddInitScript(script)
	if err != nil {
		log.Fatalf("could not add initialization script: %v", err)
	}

	log.Printf("Browser Launched, user agent: %v, Proxy: %v : %v \n", device.UserAgent, proxy.ip, proxy.pw)
	log.Println()
	return browser, nil
}

func Task(proxy ProxyStruct, wbKey string, errch chan error) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Recovered from panic: %v assuming bad proxy", err)
			time.Sleep(15 * time.Second)
			err := errors.New("recovered from panic")
			errch <- err

		}

	}()

	log.Println("Initializing playwright instance")
	PlaywrightInstance, err := playwright.Run()
	if err != nil {
		log.Panicf("could not start playwright: %v\n", err)
	}
	log.Println("Initializing browser")
	browser, err := PlaywrightInit(proxy, PlaywrightInstance)
	if err != nil {
		log.Panicf("could not launch browser: %v", err)
	}

	defer func(browser playwright.BrowserContext) {
		err := browser.Close()
		if err != nil {

		}
	}(browser)

	defer func(PlaywrightInstance *playwright.Playwright) {
		err := PlaywrightInstance.Stop()
		if err != nil {
			log.Fatalf("could not stop playwright: %v", err)
		}
	}(PlaywrightInstance)

	log.Println("Starting Task for proxy: ", proxy.ip)
	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)

	}
	if _, err = page.Goto("https://www.whatsmyip.org/"); err != nil {
		log.Panicf("could not goto: %v", err)

	}

	page2, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)

	}
	if _, err = page2.Goto("https://bot.sannysoft.com/"); err != nil {
		log.Panicf("could not goto: %v", err)
	}
	if _, err = page2.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String("botdetect.png"),
	}); err != nil {
		log.Fatalf("could not create screenshot: %v", err)
	}
	AssertErrorToNil("Failed to close Page", page.Close())
	AssertErrorToNil("Failed to close Page", page2.Close())

	time.Sleep(2 * time.Second)
	log.Println("First checks done and creenshot of sannysoft saved")

	page3, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)

	}
	if _, err = page3.Goto("https://cpr.ticketera.com/tickets/series/badbunnymostwantedtour"); err != nil {
		log.Panicf("could not goto: %v", err)
	}
	log.Println("Hit ticketera, now sleeping for 30 seconds")
	time.Sleep(30 * time.Second)
	if _, err = page3.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String("Captchaload.png"),
	}); err != nil {
		log.Fatalf("could not create screenshot: %v", err)
	}
	if page3.URL() == "https://cpr.ticketera.com/tickets/series/badbunnymostwantedtour" {
		if IdFound("#log-in-here", page3) == true {
			log.Println("Hit landing page, checking for tickets")
			//		WebhookPageLive(page3.URL(), wbKey)
			if CheckForTickets("https://cpr.ticketera.com/tickets/series/badbunnymostwantedtour/bad-bunny-most-wanted-tour-965527/bestAvailable?startDate=06-07-2024", page3) {
				WebhookTicketsLive(page3.URL(), wbKey)
			}
			time.Sleep(2 * time.Second)
			if CheckForTickets("https://cpr.ticketera.com/tickets/series/badbunnymostwantedtour/bad-bunny-most-wanted-tour-966915/bestAvailable?startDate=06-08-2024", page3) {
				WebhookTicketsLive(page3.URL(), wbKey)
			}

		} else if IdFound("#challenge-stage", page3) == true {
			webhookPageCaptchaStuck(proxy.ip, wbKey)
			log.Println("Tas Captcha stuck")

		}
	} else if page3.URL() == "https://tixtrack.queue-it.net/afterevent.aspx?c=tixtrack&e=badbunnymostwanted&t=https%3A%2F%2Fcpr.ticketera.com%2Ftickets%2Fseries%2Fbadbunnymostwantedtour&cid=en-US" {
		if TextBodyFound("The event has ended", page3) == true {
			log.Printf("queue is closed")
		} else if TextBodyFound("The event will beging", page3) == true {
			WebhookQueueUp(page3.URL(), wbKey)
			log.Println("Queue is up")
		}

	} else {
		log.Printf("somethign went wrong")
	}

	time.Sleep(3 * time.Second)
	log.Printf("finished task for proxy :%v", proxy.ip)
	AssertErrorToNil("Failed to close Page", page3.Close())
	log.Println("Task finished")
	errch <- nil

}
