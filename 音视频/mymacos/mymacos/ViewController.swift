//
//  ViewController.swift
//  mymacos
//
//  Created by WH37 on 2022/9/15.
//

import Cocoa


class ViewController: NSViewController {

    var recStatus: Bool = false
    var thread: Thread?
    let btn = NSButton.init(title: "", target: nil, action: nil)
    
    override func viewDidLoad() {
        super.viewDidLoad()
        
        // Do any additional setup after loading the view.
        self.view.setFrameSize(NSSize(width: 320, height: 240))
        btn.title = "开始录制"
        btn.bezelStyle = .rounded
        btn.setButtonType(.pushOnPushOff)
        btn.frame = NSRect(x: 320/2-60, y: 240/2-15, width: 110, height: 30)
        btn.target = self
        btn.action = #selector(myfunc)
        
        self.view.addSubview(btn)
    }
    
    @objc func myfunc() {
        self.recStatus = !self.recStatus
        
        if recStatus {
            thread = Thread.init(target: self,
                                 selector: #selector(self.recAudio),
                                 object: nil)
            thread?.start()
            self.btn.title = "停止录制"
        } else {
            set_status(0)
            self.btn.title = "开始录制"
        }

    }
    
    @objc func recAudio() {
//        print("start thread")
        rec_audio()
    }

    override var representedObject: Any? {
        didSet {
        // Upd ate the view, if already loaded.
        }
    }
}
