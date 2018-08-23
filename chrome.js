'use strict'
const puppeteer = require('puppeteer');


(async () => {
    
    const browser = await puppeteer.launch({
        // headless: false,
        args: [
            '--no-sandbox', '--disable-setuid-sandbox',
            // '--proxy-server=199.115.116.233:4000',
        ]
    });
    const page = await browser.newPage();

    try {
        var url = process.argv[2]
        await page.goto(url, {
            timeout: 30000,
            waitUntil: ["domcontentloaded"]
        });
    
        var dom = await page.content()
        console.log(dom)
    } catch (error) {
        console.log(error)
    }finally{
        await page.close()
        await browser.close()
    }
})()