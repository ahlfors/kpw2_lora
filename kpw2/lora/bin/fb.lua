local ffi = require("ffi")
local blitbuffer = require("ffi/blitbuffer")
local fb = require("ffi/framebuffer").open("/dev/fb0")
local evloop = require("ffi/eventloop")
local input = require("ffi/input")
local posix = require("ffi/posix_h")
local rfb = require("ffi/rfbclient")

local password = nil
local client = nil
local rfbFramebuffer = nil
local configfile = "config.lua"

local waitRefresh = 250
local rotateFB = 0
local reconnecting = false
local debug = false
local blitfunc = nil

local update_x1 = nil
local update_x2 = 0
local update_y1 = 0
local update_y2 = 0

local refresh_full_every_256th_pxup = 512
local refresh_full_ctr = 0

-- this is just an ID value
local TIMER_REFRESH = 10

-- constants for screen updates
local WAVEFORM_MODE_INIT      = 0x0   -- Screen goes to white (clears)
local WAVEFORM_MODE_DU        = 0x1   -- Grey->white/grey->black
local WAVEFORM_MODE_GC16      = 0x2   -- High fidelity (flashing)
local WAVEFORM_MODE_GC4       = WAVEFORM_MODE_GC16 -- For compatibility
local WAVEFORM_MODE_GC16_FAST = 0x3   -- Medium fidelity
local WAVEFORM_MODE_A2        = 0x4   -- Faster but even lower fidelity
local WAVEFORM_MODE_GL16      = 0x5   -- High fidelity from white transition
local WAVEFORM_MODE_GL16_FAST = 0x6   -- Medium fidelity from white transition
local WAVEFORM_MODE_AUTO      = 0x101

local waveform_default_fast = WAVEFORM_MODE_GC16
local waveform_default_slow = WAVEFORM_MODE_GC16

local rfbFramebuffer = nil
local rotateFB = 270
local reconnecting = false


fb.bb:rotate(rotateFB)
print(fb:getSize())

print(ffi.string(fb.finfo.id, 7))
print("************fbsize*********")
print(fb.vinfo.xres_virtual)
print(fb.vinfo.yres_virtual)
print(fb.vinfo.bits_per_pixel)
print(fb.fb_size)
print("************fbsize*********")
print(fb.vinfo.xres)
print(fb.vinfo.yres)
print(fb.finfo.smem_len)
print(fb.finfo.smem_start)
print(fb.finfo.line_length)
local str = "*****************"



local function refreshTimerFunc()
	-- not sure how this could happen but it does.
	-- TODO: find race condition
	--if not update_x1 then return end

	local x = update_x1
	local y = update_y1
	local w = update_x2 - update_x1
	local h = update_y2 - update_y1

	if debug then
		io.stdout:write(
			"eink update ", x, ",", y, " ",
			w, "x", h, "\n")
	end

	fb.bb:blitFrom(rfbFramebuffer,
		x, y, x, y, w, h, blitfunc)

	if do_refresh_full(w, h) then
		if debug then
			io.stdout:write("slow eink refresh\n")
		end
		fb:refresh(1, waveform_default_slow)
	else
		if debug then
			io.stdout:write("fast eink refresh\n")
		end
		fb:refresh(0, waveform_default_fast, x, y, w, h)
	end
	update_x1 = nil
end

-- 将图片文件的坐标数据输出到framebuffer中
local function paintPicture(bb, fpn)


    -- local f = readfile("/mnt/us/learn/kvncviewer/demo.pbm", "r")
    -- f:close()

    -- print(fpn)

    -- 以只读方式打开文件
    file = io.open(fpn, "r")

    -- 设置默认输入文件为 test.lua
    io.input(file)

    -- 输出文件第一行
    -- print(file:lines())
    local fileData
    local counter = 1
    local colorWhite = blitbuffer.Color4L(bit.bnot(255))
    local colorBlack = blitbuffer.Color4L(bit.bnot(0))
    repeat
            fileData = io.read()
            -- print("line", counter, ":", fileData)
            counter = counter + 1

            if fileData ~= nil then
                -- 计算坐标
                local x = 0
                local y = 0
                local color = 0
                local valueCounter = 0
                local delimiter = 44 -- "," 分隔符的ASCII码
                local delimiterCounter = 0
                for count = 1 , #fileData do
                    local number = fileData:byte(count)
                    if number == delimiter then
                        valueCounter = 0
                        delimiterCounter = delimiterCounter + 1
                    else
                        if delimiterCounter == 0 then
                            x = x * (valueCounter * 10)
                            x = x + number - 48
                            valueCounter = 1
                        end

                        if delimiterCounter == 1 then
                            y = y * (valueCounter * 10)
                            y = y + number - 48
                            valueCounter = 1
                        end

                        if delimiterCounter == 2 then
                            color = color * (valueCounter * 10)
                            color = color + number - 48
                            valueCounter = 1
                        end
                    end

                end

                -- print("x: ", x)
                -- print("y: ", y)
                -- print("color: ", color)

                if color == 255 then
                    bb:setPixel(x, y, colorWhite)
                end

                if color == 0 then
                    bb:setPixel(x, y, colorBlack)
                end

            end

    until fileData ==nil

    -- 关闭打开的文件
    io.close(file)

end


local function drawLine(bb, x, y)
    for m=0, y-1 do
        for n=x, x+10 do
            local c = blitbuffer.Color4L(255)
            bb:setPixel(n, m, c)
        end
        -- print(m)
    end
end

-- fb.bb = blitbuffer.fromstring(500, 200, blitbuffer.TYPE_BB8, str, nil)
-- fb.bb:invert()
-- fb.bb:invert()
-- fb.bb:invert()
-- fb.einkUpdateFunc = k51_update
-- fb:refresh(2, 1, 0, 0, 758, 1024)
-- print(fb.data)
-- fb.bb.fromstring(100, 100, 0, "hello", nil)

-- fb.einkWaitForSubmissionFunc(fb)
fb.bb.data = fb.data
repeat
	-- paintRect(x, y, w, h, value)
	-- fb.bb:paintRect(350, 300, 300, 300, 255)
    -- local x = 758/2
    -- local y = 1024

    -- drawLine(fb.bb, x, y)
    -- paintPicture(fb.bb, "/mnt/us/learn/kvncviewer/metricwhite.txt")
    paintPicture(fb.bb, "/mnt/us/extensions/lora/bin/metricwhite.txt")

    paintPicture(fb.bb, "/mnt/us/extensions/lora/bin/metric.txt")

	-- fb.bb:paintCircle(500, 300, 200, 255, 20)
    -- local function k51_update(fb, refreshtype, waveform_mode, x, y, w, h)
    -- refreshtype  -- 0 - 局刷; 1 - 全刷
	fb.einkUpdateFunc(fb, 0, WAVEFORM_MODE_GC16, 0, 0, 758, 1024)

	-- local client = rfb.rfbGetClient(8,3,4)	
	
	-- rfbFrameBuffer = blitbuffer:fromstring(200, 50, blitbuffer.TYPE_BBBGR32, client.frameBuffer)
	-- rfbFrameBuffer:invert()
	
	ffi.C.sleep(1)
until false
-- fb.bb:invert()
-- repeat
	-- local client = rfb.rfbGetClient(8,3,4)
	-- client.canHandleNewFBSize = 0
	-- client.GotFrameBufferUpdate = updateFromRFB
			
	-- rfbFramebuffer = blitbuffer.new(200, 200, blitbuffer.TYPE_BB8, client.frameBuffer)
	-- rfbFramebuffer:invert()

	-- ffi.C.sleep(2)

-- until not reconnecting

-- fb:close()

