//
//  ViewController.swift
//  mymacos
//
//  Created by WH37 on 2022/9/15.
//

import Cocoa


class ViewController: NSViewController {

    override func viewDidLoad() {
        super.viewDidLoad()
        
        // Do any additional setup after loading the view.
        self.view.setFrameSize(NSSize(width: 320, height: 240))
        
        let btn = NSButton.init(title: "button", target: nil, action: nil)
        btn.title = "hello"
        btn.frame = NSRect(x: 320/2-40, y: 240/2-15, width: 80, height: 30)
        btn.bezelStyle = .rounded
        btn.setButtonType(.pushOnPushOff)
        
        // callback
        btn.target = self
        btn.action = #selector(myfunc)
        
        self.view.addSubview(btn)
    }
    
    @objc
    func myfunc() {
        haha()
    }

    override var representedObject: Any? {
        didSet {
        // Upd ate the view, if already loaded.
        }
    }
}
