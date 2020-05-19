Pod::Spec.new do |spec|
  spec.name         = 'Moac'
  spec.version      = '{{.Version}}'
  spec.license      = { :type => 'GNU Lesser General Public License, Version 3.0' }
  spec.homepage     = 'https://github.com/filestorm/go-filestorm/moac/MoacCore'
  spec.authors      = { {{range .Contributors}}
		'{{.Name}}' => '{{.Email}}',{{end}}
	}
  spec.summary      = 'iOS MoacNode Client'
  spec.source       = { :git => 'https://github.com/filestorm/go-filestorm/moac/MoacCore.git', :commit => '{{.Commit}}' }

	spec.platform = :ios
  spec.ios.deployment_target  = '9.0'
	spec.ios.vendored_frameworks = 'Frameworks/Moac.framework'

	spec.prepare_command = <<-CMD
    curl https://gethstore.blob.core.windows.net/builds/{{.Archive}}.tar.gz | tar -xvz
    mkdir Frameworks
    mv {{.Archive}}/Moac.framework Frameworks
    rm -rf {{.Archive}}
  CMD
end
