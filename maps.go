package myGin

type handlemap map[string]HandlersChain

type methodMap struct {
	method string
	hmap   handlemap
}

type methodMaps []methodMap

func (maps methodMaps) get(method string) handlemap{
	for _, m := range maps {
		if m.method == method {
			return m.hmap
		}
	}
	return nil
}

func (h handlemap) addRoute(path string, handlers HandlersChain) {
	h[path] = handlers
}