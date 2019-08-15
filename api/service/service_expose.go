package service

import (
	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	"github.com/godbus/dbus/prop"
	"github.com/muka/go-bluetooth/bluez"
)

type ExposedDBusService interface {
	Path() dbus.ObjectPath
	Interface() string
	GetProperties() bluez.Properties
	App() *App
	DBusProperties() *DBusProperties
}

type AppService interface {
	App() *App
	Path() dbus.ObjectPath
}

func RemoveDBusService(s ExposedDBusService) error {

	err := s.App().ObjectManager().RemoveObject(s.Path())
	if err != nil {
		return err
	}

	err = s.App().exportTree()
	if err != nil {
		return err
	}

	return nil
}

func ExposeDBusService(s ExposedDBusService) error {

	conn, err := dbus.SystemBus()
	if err != nil {
		return err
	}

	err = conn.Export(s, s.Path(), s.Interface())
	if err != nil {
		return err
	}

	err = s.DBusProperties().AddProperties(s.Interface(), s.GetProperties())
	if err != nil {
		return err
	}

	s.DBusProperties().Expose(s.Path())

	node := &introspect.Node{
		Interfaces: []introspect.Interface{
			//Introspect
			introspect.IntrospectData,
			//Properties
			prop.IntrospectData,
			{
				Name:       s.Interface(),
				Methods:    introspect.Methods(s),
				Properties: s.DBusProperties().Introspection(s.Interface()),
			},
		},
	}

	// fmt.Printf("ExposeService\n\n%++v\n\n", propInterface.Introspection(s.Interface()))

	err = conn.Export(
		introspect.NewIntrospectable(node),
		s.Path(),
		"org.freedesktop.DBus.Introspectable")
	if err != nil {
		return err
	}

	return nil
}